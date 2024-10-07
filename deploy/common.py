import os
import xml.etree.ElementTree as ET

def project_path():
    return os.popen("git rev-parse --show-toplevel --show-superproject-working-tree").read().strip()

def addresses_from_manifest(manifest_file: str) -> "list[str]":
    tree = ET.parse(manifest_file)
    root = tree.getroot()
    addresses = []
    for child in root:
        # print(child.tag)
        if child.tag.endswith("node"):
            component_id = child.attrib["component_id"]
            # print(component_id)
            node_name = component_id.split("+")[-1]
            location = component_id.split("+")[1]
            address = f'{node_name}.{location}'
            # print(address)
            addresses.append(address)
    return addresses

PROJECT_PATH = project_path()
SERVERS = addresses_from_manifest(f'{PROJECT_PATH}/manifest.xml')

print("SERVERS => " + SERVERS[0])
print("PROJECT_PATH => " + PROJECT_PATH)
