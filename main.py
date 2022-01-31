import os
import sys
import requests


def print_exit(content, exit_status=0):
    print(content)
    input("Press ENTER to exit...")
    sys.exit(exit_status)


if getattr(sys, 'frozen', False):
    script_path = os.path.dirname(sys.executable)
elif __file__:
    script_path = os.path.dirname(__file__)
else:
    print_exit("Unable to determine the scripts path.", 1)

surfshark_tunnel_path = "C:\\ProgramData\\Surfshark\\WireguardConfigs\\SurfsharkWireGuard.conf"

if not os.path.exists(surfshark_tunnel_path):
    print_exit("\nSurfsharkWireGuard.conf\n doesn't exist. Run Surfshark and connect to a server using the WireGuard protocol.", 1)

with open(surfshark_tunnel_path, "r") as surfshark_tunnel:
    private_key = None
    for line in surfshark_tunnel.read().splitlines():
        if line.startswith("PrivateKey = "):
            private_key = line.removeprefix("PrivateKey = ")
            break

if not private_key:
    print_exit("Unable to parse private key.", 1)

print(f"Your private key is \"{private_key}\".")

surfshark_api_endpoints = [
    "https://api.surfshark.com/v4/server/clusters/generic?countryCode=",
    "https://api.surfshark.com/v4/server/clusters/double?countryCode=",
    "https://api.surfshark.com/v4/server/clusters/static?countryCode=",
    "https://api.surfshark.com/v4/server/clusters/obfuscated?countryCode="
]

surfshark_locations = []
for endpoint in surfshark_api_endpoints:
    response = requests.get(endpoint)
    if response.status_code != 200:
        print_exit(f"\"{endpoint}\" API call status code \"{response.status_code}\".", 1)
    for location in response.json():
        surfshark_locations.append(location)

if not surfshark_locations:
    print_exit("No locations found", 1)

tunnel_directory = "wireguard-tunnels"
tunnel_path = f"{script_path}\\{tunnel_directory}"
if not os.path.exists(tunnel_path):
    os.makedirs(tunnel_path)

for location in surfshark_locations:
    location_name = location["connectionName"]

    if not location["pubKey"]:
        print(f"""\"{location_name}\" has no public key""")
        continue

    with open(f"""{tunnel_path}\\{location_name}.conf""", "w") as tunnel:
        tunnel.write(f"""[Interface]\nPrivateKey = {private_key}\nAddress = 10.14.0.2/16\nDNS = 162.252.172.57, 149.154.159.92\n[Peer]\nPublicKey = {location["pubKey"]}\nAllowedIps = 0.0.0.0/0\nEndpoint = {location_name}:51820\n[Peer]\nPublicKey = o07k/2dsaQkLLSR0dCI/FUd3FLik/F/HBBcOGUkNQGo=\nAllowedIPs = 172.16.0.36/32\nEndpoint = 92.249.38.1:51820\n""")

print_exit(f"WireGuard tunnel generation complete. Check the \"{tunnel_directory}\" folder.")
