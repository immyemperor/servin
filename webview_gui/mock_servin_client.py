"""
Demo Servin Client for Windows
This is a mock implementation for demonstration purposes when the actual servin binary is not available for Windows
"""

import json
import os
import time
from datetime import datetime
from typing import List, Dict, Any, Optional

class ServinError(Exception):
    """Exception raised for servin command errors"""
    pass

class MockServinClient:
    """Mock Servin client for demonstration purposes"""
    
    def __init__(self, servin_path: Optional[str] = None):
        """Initialize the mock ServinClient"""
        self.servin_path = servin_path or "servin"
        print("Mock Servin Client initialized (demo mode)")
        
        # Mock data
        self._containers = [
            {
                'id': 'abc123456789',
                'name': 'demo-nginx',
                'image': 'nginx:latest',
                'status': 'running',
                'state': 'running',
                'created': datetime.now().isoformat(),
                'ports': ['8080:80'],
                'networks': ['bridge']
            },
            {
                'id': 'def987654321',
                'name': 'demo-ubuntu',
                'image': 'ubuntu:20.04',
                'status': 'stopped',
                'state': 'stopped',
                'created': datetime.now().isoformat(),
                'ports': [],
                'networks': ['bridge']
            }
        ]
        
        self._images = [
            {
                'id': 'sha256:abc123',
                'repository': 'nginx',
                'tag': 'latest',
                'created': datetime.now().isoformat(),
                'size': 133000000,
                'virtual_size': 133000000
            },
            {
                'id': 'sha256:def456',
                'repository': 'ubuntu',
                'tag': '20.04',
                'created': datetime.now().isoformat(),
                'size': 72800000,
                'virtual_size': 72800000
            }
        ]
        
        self._volumes = [
            {
                'name': 'demo-data',
                'driver': 'local',
                'mountpoint': '/var/lib/servin/volumes/demo-data',
                'created': datetime.now().isoformat(),
                'scope': 'local'
            },
            {
                'name': 'demo-logs',
                'driver': 'local',
                'mountpoint': '/var/lib/servin/volumes/demo-logs',
                'created': datetime.now().isoformat(),
                'scope': 'local'
            }
        ]
    
    def ping(self) -> bool:
        """Test if servin is working"""
        return True
    
    # Container Management Methods
    
    def list_containers(self, all_containers: bool = True) -> List[Dict[str, Any]]:
        """List containers"""
        return self._containers.copy()
    
    def get_container(self, container_id: str) -> Dict[str, Any]:
        """Get detailed information about a specific container"""
        for container in self._containers:
            if container['id'].startswith(container_id) or container['name'] == container_id:
                return container
        raise ServinError(f"Container not found: {container_id}")
    
    def start_container(self, container_id: str) -> bool:
        """Start a container"""
        for container in self._containers:
            if container['id'].startswith(container_id) or container['name'] == container_id:
                if container['status'] == 'stopped':
                    container['status'] = 'running'
                    container['state'] = 'running'
                    return True
                else:
                    raise ServinError("Container is already running")
        raise ServinError(f"Container not found: {container_id}")
    
    def stop_container(self, container_id: str) -> bool:
        """Stop a container"""
        for container in self._containers:
            if container['id'].startswith(container_id) or container['name'] == container_id:
                if container['status'] == 'running':
                    container['status'] = 'stopped'
                    container['state'] = 'stopped'
                    return True
                else:
                    raise ServinError("Container is already stopped")
        raise ServinError(f"Container not found: {container_id}")
    
    def restart_container(self, container_id: str) -> bool:
        """Restart a container"""
        self.stop_container(container_id)
        time.sleep(1)  # Simulate restart delay
        self.start_container(container_id)
        return True
    
    def remove_container(self, container_id: str, force: bool = False) -> bool:
        """Remove a container"""
        for i, container in enumerate(self._containers):
            if container['id'].startswith(container_id) or container['name'] == container_id:
                if container['status'] == 'running' and not force:
                    raise ServinError("Cannot remove running container. Use force=True or stop it first.")
                del self._containers[i]
                return True
        raise ServinError(f"Container not found: {container_id}")
    
    # Image Management Methods
    
    def list_images(self) -> List[Dict[str, Any]]:
        """List images"""
        return self._images.copy()
    
    def pull_image(self, image_name: str) -> bool:
        """Pull an image (mock)"""
        # Add a new mock image
        new_image = {
            'id': f'sha256:{hash(image_name + str(time.time()))}',
            'repository': image_name.split(':')[0] if ':' in image_name else image_name,
            'tag': image_name.split(':')[1] if ':' in image_name else 'latest',
            'created': datetime.now().isoformat(),
            'size': 50000000,
            'virtual_size': 50000000
        }
        self._images.append(new_image)
        return True
    
    def remove_image(self, image_id: str, force: bool = False) -> bool:
        """Remove an image"""
        for i, image in enumerate(self._images):
            if image['id'].startswith(image_id) or f"{image['repository']}:{image['tag']}" == image_id:
                del self._images[i]
                return True
        raise ServinError(f"Image not found: {image_id}")
    
    def import_image(self, tarball_path: str, image_name: str) -> bool:
        """Import an image from tarball (mock)"""
        if not os.path.exists(tarball_path):
            raise ServinError(f"Tarball not found: {tarball_path}")
        
        new_image = {
            'id': f'sha256:{hash(image_name + str(time.time()))}',
            'repository': image_name.split(':')[0] if ':' in image_name else image_name,
            'tag': image_name.split(':')[1] if ':' in image_name else 'latest',
            'created': datetime.now().isoformat(),
            'size': 75000000,
            'virtual_size': 75000000
        }
        self._images.append(new_image)
        return True
    
    # Volume Management Methods
    
    def list_volumes(self) -> List[Dict[str, Any]]:
        """List volumes"""
        return self._volumes.copy()
    
    def create_volume(self, name: str, driver: str = "local") -> bool:
        """Create a volume"""
        # Check if volume already exists
        for volume in self._volumes:
            if volume['name'] == name:
                raise ServinError(f"Volume already exists: {name}")
        
        new_volume = {
            'name': name,
            'driver': driver,
            'mountpoint': f'/var/lib/servin/volumes/{name}',
            'created': datetime.now().isoformat(),
            'scope': 'local'
        }
        self._volumes.append(new_volume)
        return True
    
    def remove_volume(self, volume_name: str, force: bool = False) -> bool:
        """Remove a volume"""
        for i, volume in enumerate(self._volumes):
            if volume['name'] == volume_name:
                del self._volumes[i]
                return True
        raise ServinError(f"Volume not found: {volume_name}")
    
    # System Information Methods
    
    def info(self) -> Dict[str, Any]:
        """Get system information"""
        running_containers = len([c for c in self._containers if c['status'] == 'running'])
        stopped_containers = len([c for c in self._containers if c['status'] == 'stopped'])
        
        return {
            'containers': len(self._containers),
            'containers_running': running_containers,
            'containers_paused': 0,
            'containers_stopped': stopped_containers,
            'images': len(self._images),
            'server_version': 'Servin 1.0 (Demo Mode)',
            'operating_system': 'Linux',
            'architecture': 'x86_64',
            'memory_total': 8000000000,  # 8GB
            'cpu_count': 4
        }

    def inspect_container(self, container_id: str) -> Dict[str, Any]:
        """Get detailed information about a container"""
        container = self.get_container(container_id)
        
        # Add detailed inspect information
        container.update({
            'config': {
                'image': container['image'],
                'cmd': ['/bin/sh', '-c', 'while true; do sleep 30; done;'],
                'env': [
                    'PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin',
                    'HOME=/root',
                    'HOSTNAME=' + container_id[:12]
                ],
                'working_dir': '/',
                'hostname': container_id[:12]
            },
            'network_settings': {
                'ip_address': '172.17.0.2',
                'gateway': '172.17.0.1',
                'mac_address': '02:42:ac:11:00:02',
                'ports': [
                    {'container_port': '80', 'host_port': '8080', 'protocol': 'tcp'}
                ] if container['ports'] else []
            },
            'mounts': [
                {'source': '/host/data', 'destination': '/data', 'mode': 'rw', 'type': 'bind'}
            ],
            'state': {
                'status': container['status'],
                'running': container['status'] == 'running',
                'pid': 1234 if container['status'] == 'running' else 0,
                'exit_code': 0 if container['status'] == 'running' else 1,
                'started_at': container['created'],
                'finished_at': '' if container['status'] == 'running' else container['created']
            }
        })
        
        return container

    def get_logs(self, container_id: str, follow: bool = False, tail: int = 100) -> str:
        """Get container logs"""
        container = self.get_container(container_id)
        
        # Generate mock logs
        logs = [
            f"2024-09-15 12:00:0{i} [INFO] Container {container['name']} starting..."
            for i in range(5)
        ]
        logs.extend([
            f"2024-09-15 12:0{i}:00 [INFO] Processing request {i}"
            for i in range(1, min(tail, 10))
        ])
        
        if container['status'] == 'running':
            logs.append("2024-09-15 12:10:00 [INFO] Container is healthy and running")
        else:
            logs.append("2024-09-15 12:10:00 [INFO] Container stopped")
        
        return '\n'.join(logs)

    def list_files(self, container_id: str, path: str = '/') -> List[Dict[str, Any]]:
        """List files in container filesystem"""
        # Mock filesystem structure with proper path handling
        if path == '/':
            return [
                {'name': 'bin', 'type': 'directory', 'size': 4096, 'permissions': 'drwxr-xr-x', 'is_dir': True, 'path': '/bin'},
                {'name': 'etc', 'type': 'directory', 'size': 4096, 'permissions': 'drwxr-xr-x', 'is_dir': True, 'path': '/etc'},
                {'name': 'home', 'type': 'directory', 'size': 4096, 'permissions': 'drwxr-xr-x', 'is_dir': True, 'path': '/home'},
                {'name': 'opt', 'type': 'directory', 'size': 4096, 'permissions': 'drwxr-xr-x', 'is_dir': True, 'path': '/opt'},
                {'name': 'tmp', 'type': 'directory', 'size': 4096, 'permissions': 'drwxrwxrwx', 'is_dir': True, 'path': '/tmp'},
                {'name': 'usr', 'type': 'directory', 'size': 4096, 'permissions': 'drwxr-xr-x', 'is_dir': True, 'path': '/usr'},
                {'name': 'var', 'type': 'directory', 'size': 4096, 'permissions': 'drwxr-xr-x', 'is_dir': True, 'path': '/var'},
                {'name': 'app.log', 'type': 'file', 'size': 1024, 'permissions': '-rw-r--r--', 'is_dir': False, 'path': '/app.log'}
            ]
        elif path == '/etc':
            return [
                {'name': 'passwd', 'type': 'file', 'size': 2048, 'permissions': '-rw-r--r--', 'is_dir': False, 'path': '/etc/passwd'},
                {'name': 'hosts', 'type': 'file', 'size': 256, 'permissions': '-rw-r--r--', 'is_dir': False, 'path': '/etc/hosts'},
                {'name': 'hostname', 'type': 'file', 'size': 12, 'permissions': '-rw-r--r--', 'is_dir': False, 'path': '/etc/hostname'},
                {'name': 'ssl', 'type': 'directory', 'size': 4096, 'permissions': 'drwxr-xr-x', 'is_dir': True, 'path': '/etc/ssl'}
            ]
        elif path == '/home':
            return [
                {'name': 'user', 'type': 'directory', 'size': 4096, 'permissions': 'drwxr-xr-x', 'is_dir': True, 'path': '/home/user'},
                {'name': 'admin', 'type': 'directory', 'size': 4096, 'permissions': 'drwxr-xr-x', 'is_dir': True, 'path': '/home/admin'}
            ]
        elif path == '/var':
            return [
                {'name': 'log', 'type': 'directory', 'size': 4096, 'permissions': 'drwxr-xr-x', 'is_dir': True, 'path': '/var/log'},
                {'name': 'www', 'type': 'directory', 'size': 4096, 'permissions': 'drwxr-xr-x', 'is_dir': True, 'path': '/var/www'},
                {'name': 'lib', 'type': 'directory', 'size': 4096, 'permissions': 'drwxr-xr-x', 'is_dir': True, 'path': '/var/lib'}
            ]
        else:
            # For any other path, return some sample files
            return [
                {'name': 'file1.txt', 'type': 'file', 'size': 512, 'permissions': '-rw-r--r--', 'is_dir': False, 'path': f'{path}/file1.txt'},
                {'name': 'file2.txt', 'type': 'file', 'size': 1024, 'permissions': '-rw-r--r--', 'is_dir': False, 'path': f'{path}/file2.txt'},
                {'name': 'config.conf', 'type': 'file', 'size': 2048, 'permissions': '-rw-r--r--', 'is_dir': False, 'path': f'{path}/config.conf'}
            ]

    def exec_command(self, container_id: str, command: str) -> str:
        """Execute command in container"""
        container = self.get_container(container_id)
        
        if container['status'] != 'running':
            raise ServinError(f"Container {container_id} is not running")
        
        # Mock command responses
        if command.startswith('ls'):
            return "app.log  bin  etc  home  opt  tmp  usr  var"
        elif command == 'env':
            return "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\nHOME=/root\nHOSTNAME=" + container_id[:12]
        elif command == 'pwd':
            return "/"
        elif command == 'whoami':
            return "root"
        else:
            return f"Mock output for command: {command}"

    def get_environment(self, container_id: str) -> List[Dict[str, str]]:
        """Get container environment variables"""
        container = self.get_container(container_id)
        
        return [
            {'key': 'PATH', 'value': '/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin'},
            {'key': 'HOME', 'value': '/root'},
            {'key': 'HOSTNAME', 'value': container_id[:12]},
            {'key': 'TERM', 'value': 'xterm'},
            {'key': 'CONTAINER_NAME', 'value': container['name']},
            {'key': 'CONTAINER_IMAGE', 'value': container['image']}
        ]

# Use the mock client when the real one is not available
try:
    from servin_client import ServinClient, ServinError
    # Try to create a real client
    test_client = ServinClient()
    # If successful, use the real client
    print("Using real Servin client")
except:
    # If failed, use the mock client
    print("Using mock Servin client for demo")
    ServinClient = MockServinClient
