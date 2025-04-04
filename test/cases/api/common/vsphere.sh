#!/usr/bin/bash

source /usr/libexec/tests/osbuild-composer/shared_lib.sh

function installClientVSphere() {
    if ! hash govc; then
        ARCH="$(uname -m)"
        if [ "$ARCH" = "aarch64" ]; then
            ARCH="arm64"
        fi
        greenprint "Installing govc"
        pushd "${WORKDIR}" || exit 1
        curl -Ls --retry 5 --output govc.tar.gz \
            "https://github.com/vmware/govmomi/releases/download/v0.29.0/govc_Linux_$ARCH.tar.gz"
        tar -xvf govc.tar.gz
        GOVC_CMD="${WORKDIR}/govc"
        chmod +x "${GOVC_CMD}"
        popd || exit 1
    else
        echo "Using pre-installed 'govc' from the system"
        GOVC_CMD="govc"
    fi

    $GOVC_CMD version
}

function checkEnvVSphere() {
    printenv VC8_GOVMOMI_USERNAME VC8_GOVMOMI_PASSWORD VC8_GOVMOMI_URL VC8_GOVMOMI_CLUSTER VC8_GOVMOMI_DATACENTER VC8_GOVMOMI_DATASTORE VC8_GOVMOMI_FOLDER VC8_GOVMOMI_NETWORK  > /dev/null
}

# Create a cloud-int user-data file
#
# Returns:
#   - path to the user-data file
#
# Arguments:
#   $1 - default username
#   $2 - path to the SSH public key to set as authorized for the user
function createCIUserdata() {
    local _user="$1"
    local _ssh_pubkey_path="$2"

    local _ci_userdata_dir
    _ci_userdata_dir="$(mktemp -d -p "${WORKDIR}")"
    local _ci_userdata_path="${_ci_userdata_dir}/user-data"

    cat > "${_ci_userdata_path}" <<EOF
#cloud-config
users:
    - name: "${_user}"
      sudo: "ALL=(ALL) NOPASSWD:ALL"
      ssh_authorized_keys:
          - "$(cat "${_ssh_pubkey_path}")"
EOF

    echo "${_ci_userdata_path}"
}

# Create a cloud-int meta-data file
#
# Returns:
#   - path to the meta-data file
#
# Arguments:
#   $1 - VM name
function createCIMetadata() {
    local _vm_name="$1"

    local _ci_metadata_dir
    _ci_metadata_dir="$(mktemp -d -p "${WORKDIR}")"
    local _ci_metadata_path="${_ci_metadata_dir}/meta-data"

    cat > "${_ci_metadata_path}" <<EOF
instance-id: ${_vm_name}
local-hostname: ${_vm_name}
EOF

    echo "${_ci_metadata_path}"
}

# Compress and encode the provided file path for cloud-init
#
# Returns:
#   - base64 encoded gzip compressed file as a string
#
# Arguments:
#   $1 - path to the file to encode
function fileToCloudInitEncodedGzipB64() {
    local _file_path="$1"
    gzip -c "${_file_path}" | base64 -w 0
}

# Create an ISO with the provided cloud-init user-data file
#
# Returns:
#   - path to the created ISO file
#
# Arguments:
#   $1 - path to the cloud-init user-data file
#   $2 - path to the cloud-init meta-data file
function createCIUserdataISO() {
    local _ci_userdata_path="$1"
    local _ci_metadata_path="$2"

    local _iso_path
    _iso_path="$(mktemp -p "${WORKDIR}" --suffix .iso)"
    mkisofs \
        -input-charset "utf-8" \
        -output "${_iso_path}" \
        -volid "cidata" \
        -joliet \
        -rock \
        -quiet \
        -graft-points \
        "${_ci_userdata_path}" \
        "${_ci_metadata_path}"

    echo "${_iso_path}"
}

# Verify VMDK image in VSphere
function verifyInVSphere() {
    local _filename="$1"
    greenprint "Verifying VMDK image: ${_filename}"

    VSPHERE_VM_NAME="osbuild-composer-vm-${TEST_ID}"
    VSPHERE_IMAGE_NAME="${VSPHERE_VM_NAME}.vmdk"
    mv "${_filename}" "${WORKDIR}/${VSPHERE_IMAGE_NAME}"

    # import the built VMDK image to VSphere
    # import.vmdk seems to be creating the provided directory and
    # if one with this name exists, it appends "_<number>" to the name
    greenprint "💿 ⬆️ Importing the converted VMDK image to VSphere"
    $GOVC_CMD import.vmdk \
        -u "${VC8_GOVMOMI_USERNAME}:${VC8_GOVMOMI_PASSWORD}@${VC8_GOVMOMI_URL}" \
        -k=true \
        -dc="${VC8_GOVMOMI_DATACENTER}" \
        -pool="${VC8_GOVMOMI_CLUSTER}"/Resources \
        -ds="${VC8_GOVMOMI_DATASTORE}" \
        "${WORKDIR}/${VSPHERE_IMAGE_NAME}" \
        "${VSPHERE_VM_NAME}"

    # create the VM, but don't start it
    greenprint "🖥️ Creating VM in VSphere"
    $GOVC_CMD vm.create \
        -u "${VC8_GOVMOMI_USERNAME}:${VC8_GOVMOMI_PASSWORD}@${VC8_GOVMOMI_URL}" \
        -k=true \
        -dc="${VC8_GOVMOMI_DATACENTER}" \
        -pool="${VC8_GOVMOMI_CLUSTER}"/Resources \
        -ds="${VC8_GOVMOMI_DATASTORE}" \
        -folder="${VC8_GOVMOMI_FOLDER}" \
        -net="${VC8_GOVMOMI_NETWORK}" \
        -net.adapter=vmxnet3 \
        -m=4096 -c=2 -g=rhel8_64Guest -on=true -firmware=efi \
        -disk="${VSPHERE_VM_NAME}/${VSPHERE_IMAGE_NAME}" \
        -disk.controller=scsi \
        -on=false \
        "${VSPHERE_VM_NAME}"

    # Create SSH keys to use
    local _vsphere_ssh_key="${WORKDIR}/vsphere_ssh_key"
    ssh-keygen -t rsa-sha2-512 -f "${_vsphere_ssh_key}" -C "${SSH_USER}" -N ""

    # Set cloud-init data for the VM
    local _ci_userdata_path
    _ci_userdata_path="$(createCIUserdata "${SSH_USER}" "${_vsphere_ssh_key}.pub")"
    local _ci_userdata_encoded
    _ci_userdata_encoded="$(fileToCloudInitEncodedGzipB64 "${_ci_userdata_path}")"

    local _ci_metadata_path
    _ci_metadata_path="$(createCIMetadata "${VSPHERE_VM_NAME}")"
    local _ci_metadata_encoded
    _ci_metadata_encoded="$(fileToCloudInitEncodedGzipB64 "${_ci_metadata_path}")"

    # configure the VM to use the cloud-init data
    greenprint "💿 Configuring the VM to use the cloud-init data"
    $GOVC_CMD vm.change \
        -u "${VC8_GOVMOMI_USERNAME}:${VC8_GOVMOMI_PASSWORD}@${VC8_GOVMOMI_URL}" \
        -k=true \
        -dc="${VC8_GOVMOMI_DATACENTER}" \
        -vm "${VSPHERE_VM_NAME}" \
        -e "guestinfo.userdata=${_ci_userdata_encoded}" \
        -e "guestinfo.userdata.encoding=gzip+base64" \
        -e "guestinfo.metadata=${_ci_metadata_encoded}" \
        -e "guestinfo.metadata.encoding=gzip+base64"

    # tagging vm as testing object
    $GOVC_CMD tags.attach \
        -u "${VC8_GOVMOMI_USERNAME}":"${VC8_GOVMOMI_PASSWORD}"@"${VC8_GOVMOMI_URL}" \
        -k=true \
        -c "osbuild-composer testing" gitlab-ci-test \
        "/${VC8_GOVMOMI_DATACENTER}/vm/${VC8_GOVMOMI_FOLDER}/${VSPHERE_VM_NAME}"

    # start the VM
    greenprint "🔌 Powering up the VSphere VM"
    $GOVC_CMD vm.power \
        -u "${VC8_GOVMOMI_USERNAME}:${VC8_GOVMOMI_PASSWORD}@${VC8_GOVMOMI_URL}" \
        -k=true \
        -dc="${VC8_GOVMOMI_DATACENTER}" \
        -on "${VSPHERE_VM_NAME}"

    HOST=$($GOVC_CMD vm.ip \
        -u "${VC8_GOVMOMI_USERNAME}:${VC8_GOVMOMI_PASSWORD}@${VC8_GOVMOMI_URL}" \
        -k=true \
        -v4=true \
        -dc="${VC8_GOVMOMI_DATACENTER}" \
        "${VSPHERE_VM_NAME}")
    greenprint "⏱ Waiting for the VSphere VM to respond to ssh"
    _instanceWaitSSH "${HOST}"

    _ssh="ssh -oStrictHostKeyChecking=no -i ${_vsphere_ssh_key} $SSH_USER@$HOST"
    _instanceCheck "${_ssh}"

    greenprint "✅ Successfully verified VSphere image with cloud-init"
}

function cleanupVSphere() {
    # since this function can be called at any time, ensure that we don't expand unbound variables
    GOVC_CMD="${GOVC_CMD:-}"
    VSPHERE_VM_NAME="${VSPHERE_VM_NAME:-}"
    VSPHERE_CIDATA_ISO_PATH="${VSPHERE_CIDATA_ISO_PATH:-}"

    greenprint "🧹 Cleaning up the VSphere VM"
    $GOVC_CMD vm.destroy \
        -u "${VC8_GOVMOMI_USERNAME}:${VC8_GOVMOMI_PASSWORD}@${VC8_GOVMOMI_URL}" \
        -k=true \
        -dc="${VC8_GOVMOMI_DATACENTER}" \
        "${VSPHERE_VM_NAME}"

    greenprint "🧹 Cleaning up the VSphere Datastore"
    $GOVC_CMD datastore.rm \
        -u "${VC8_GOVMOMI_USERNAME}:${VC8_GOVMOMI_PASSWORD}@${VC8_GOVMOMI_URL}" \
        -k=true \
        -dc="${VC8_GOVMOMI_DATACENTER}" \
        -ds="${VC8_GOVMOMI_DATASTORE}" \
        -f \
        "${VSPHERE_CIDATA_ISO_PATH}"

    $GOVC_CMD datastore.rm \
        -u "${VC8_GOVMOMI_USERNAME}:${VC8_GOVMOMI_PASSWORD}@${VC8_GOVMOMI_URL}" \
        -k=true \
        -dc="${VC8_GOVMOMI_DATACENTER}" \
        -ds="${VC8_GOVMOMI_DATASTORE}" \
        -f \
        "${VSPHERE_VM_NAME}"
}
