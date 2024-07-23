import nmap
import time
from tqdm import tqdm
from concurrent.futures import ThreadPoolExecutor

print(r"""
          ╫▓▓▓▓▓▓▓▓▓▓▓▓╩╩╨╨╩╩▀▓▓▓▓▓▀▓▓▓▓▓▓▌
          ╫▓▓▓▓▓▓▓▓▒   ...      ```  ²█▓▓▓▌
          ╫▓▓▓▓▓▓▓▓   ▒   `"═µ        ╫▓▓▓▌
          ╫▓▓▓▓▓▓▓▓╦µ `u      `ªu  .╗▓▓▓▓▓▌
          ╫▓▓▓▓▓▓▓▌░,░.░         ªª`╨█▓▓▓▓▌
          ╫▓▓▓▓▓▓▓Ñ░╩░░`▒N╥,.        ╫▓▓▓▓▌
          ╫╫▓▓▓▓█▓▓µ▒ `,.ªN░▒``"▒▒▒▒▓▓▓▓▓▓Ñ
          ╫▓▓▓▓█████H )░▒░  ░░   ▒" j███▓▓▌
          ╫▓▓▓█████▌   j▒"▒µ,..  ▒░▒j██▓▓▓▌
          ╫▓▓▓█████▒µ  `Ñ▒ ``"╨▒µ▒ ,▓███▓▓▌
          ╫▓▓▓████▌ 1,  ▒░`░═»»µ▒ ╥▓▓███▓▓▌
          ╫▓█████▌   `▒  ▒░,░.░░ ▄█▓▓▓█▓▓▓▌
          ╫██████      1µ `"░░" ▓███████▓▓▌
          ╫█▀▀" ªµ      `▒     ▒"▀▀█████▓▓▌
                 `ªu.    `░,.,¿`    `╨▀▀▓▓▌
                    `ªu,      ▒          ``
                        `ªuµ,¿╛      ~~═~░""")
print("j0kah an IP-scanner made for RECON! by b0urn3\n")

scanner = nmap.PortScanner()

target = input("What IP/domain would you like to scan? \n> ")

print(f"\nI got {target} as your target to scan, \nlet's go!")
response = input("""\nEnter the type of scan you want to run:
                i. SYN-ACK Scan
                ii. UDP Scan
                iii. AnonScan
                iv. Regular Scan
                v. OS Detection
                vi. Multiple IP inputs
                vii. Ping Scan
                viii. Comprehensive Scan\n> """)

def progress_indicator(duration):
    for _ in tqdm(range(duration), desc="Scanning", unit="s"):
        time.sleep(1)

def perform_scan(target, scan_type, args='', duration=30):
    try:
        progress_indicator(duration)
        scanner.scan(target, arguments=args)
        
        print("\nScan Results:")
        print(f"Target: {target}")
        print(f"Scan Type: {scan_type}")
        print("Target Status:", scanner[target].state())

        for protocol in scanner[target].all_protocols():
            ports = scanner[target][protocol].keys()
            open_ports = [port for port in ports if scanner[target][protocol][port]['state'] == 'open']
            if open_ports:
                print(f"Open {protocol.upper()} Ports: {open_ports}")
            else:
                print(f"No open {protocol.upper()} ports found.")
        
        if 'osclass' in scanner[target]:
            print("\nOS Details:")
            for osclass in scanner[target]['osclass']:
                print(f"OS: {osclass['osfamily']} {osclass['osgen']} - Accuracy: {osclass['accuracy']}%")
        
        print("\nDetailed Protocols:")
        print(scanner[target].all_protocols())

    except Exception as e:
        print(f"An error occurred: {e}")

def parallel_scan(targets, scan_type, args, duration):
    with ThreadPoolExecutor(max_workers=4) as executor:
        futures = [executor.submit(perform_scan, target, scan_type, args, duration) for target in targets]
        for future in futures:
            future.result()

scan_options = {
    "i": ("-sS -T4", 30),  # SYN-ACK Scan with higher timing template
    "ii": ("-sU -T4", 30),  # UDP Scan with higher timing template
    "iii": ("-sS -sU -D 10.0.0.1,10.0.0.2,10.0.0.3,192.168.1.11,192.168.1.12,192.168.1.13,192.168.1.14,192.168.1.15 -f --randomize-hosts -T4", 60),  # AnonScan with higher timing template
    "iv": ("-T4", 20),  # Regular Scan with higher timing template
    "v": ("-O -T4", 40),  # OS Detection with higher timing template
    "vi": ("-T4", 25),  # Multiple IP inputs with higher timing template
    "vii": ("-sn -T4", 15),  # Ping Scan with higher timing template
    "viii": ("-sS -sU -O -A -T4 --script vuln", 60),  # Comprehensive Scan with higher timing template and vulnerability scripts
}

if response in scan_options:
    scan_args, scan_duration = scan_options[response]
    if response == "vi":
        target_list = target.split(',')
        parallel_scan(target_list, response, scan_args, scan_duration)
    else:
        perform_scan(target, response, scan_args, scan_duration)
else:
    print("Invalid option selected!")
