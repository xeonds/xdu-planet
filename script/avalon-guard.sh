#!/bin/bash

API_URL=https://planet.iris.al/api/v1/admin/comment
API_TOKEN=

# Function to fetch audit comments
fetch_audit_comments() {
    response=$(curl -s "${API_URL}/audit" -H "Authorization: ${API_TOKEN}")
    echo "$response" | jq -r '.[] | "ID: \(.ID) Content: \(.content)"'
}

# Function to update comment status
update_comment_status() {
    local status=$1
    shift
    local ids=("$@")

    for id in "${ids[@]}"; do
        curl -s -X POST "${API_URL}/audit/${id}" -H "application/json" -d "{\"status\": \"${status}\"}" -H "Authorization: ${API_TOKEN}" \
        && echo "Updated comment ID ${id} to status ${status}"
    done
}

# Display usage
usage() {
    echo "Usage: $0 [-f] [-u status id1 id2 ...]"
    echo "  -f                Fetch all audit (pending) comments"
    echo "  -u status ids...  Update comment statuses. Status can be 'ok', 'block', 'audit' or 'delete'"
    exit 1
}

# Parse command line arguments
while getopts "fu:" opt; do
    case ${opt} in
        f)
            fetch_audit_comments
            ;;
        u)
            shift $((OPTIND-2))
            # if [[ -z "$1" || "$1" != "ok" && "$1" != "block" ]]; then
            #     usage
            # fi
            status=$1
            shift
            update_comment_status $status "$@"
            ;;
        *)
            usage
            ;;
    esac
done

if [ $OPTIND -eq 1 ]; then
    usage
fi
