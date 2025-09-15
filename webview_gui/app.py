"""
Servin Runtime
Flask API server for managing Servin containers, images, and volumes
"""

import os
import sys
import threading
import time
from flask import Flask, jsonify, request, render_template, send_from_directory
from flask_cors import CORS
from servin_client import ServinClient, ServinError

app = Flask(__name__)
CORS(app)  # Enable CORS for all routes

# Initialize Servin client
try:
    from servin_client import ServinClient, ServinError
    servin_client = ServinClient()
    # Test connection
    servin_client.ping()
    print("Successfully connected to Servin runtime")
except (ServinError, Exception) as e:
    print(f"Failed to connect to real Servin runtime: {e}")
    try:
        # Try to use mock client for demo
        from mock_servin_client import ServinClient, ServinError
        servin_client = ServinClient()
        print("Using mock Servin client for demonstration")
    except Exception as e2:
        print(f"Failed to initialize mock client: {e2}")
        servin_client = None

@app.route('/')
def index():
    """Serve the main HTML page"""
    import time
    timestamp = int(time.time())
    return render_template('index.html', timestamp=timestamp)

@app.route('/static/<path:filename>')
def static_files(filename):
    """Serve static files"""
    return send_from_directory('static', filename)

# Container Management APIs
@app.route('/api/containers', methods=['GET'])
def get_containers():
    """Get list of all containers"""
    if not servin_client:
        return jsonify({'error': 'Servin runtime not available'}), 500
    
    try:
        containers = servin_client.list_containers()
        return jsonify(containers)
    except ServinError as e:
        return jsonify({'error': str(e)}), 500

@app.route('/api/containers/<container_id>/start', methods=['POST'])
def start_container(container_id):
    """Start a container"""
    if not servin_client:
        return jsonify({'error': 'Servin runtime not available'}), 500
    
    try:
        servin_client.start_container(container_id)
        return jsonify({'success': True, 'message': f'Container {container_id} started'})
    except ServinError as e:
        return jsonify({'error': str(e)}), 500

@app.route('/api/containers/<container_id>/stop', methods=['POST'])
def stop_container(container_id):
    """Stop a container"""
    if not servin_client:
        return jsonify({'error': 'Servin runtime not available'}), 500
    
    try:
        servin_client.stop_container(container_id)
        return jsonify({'success': True, 'message': f'Container {container_id} stopped'})
    except ServinError as e:
        return jsonify({'error': str(e)}), 500

@app.route('/api/containers/<container_id>/restart', methods=['POST'])
def restart_container(container_id):
    """Restart a container"""
    if not servin_client:
        return jsonify({'error': 'Servin runtime not available'}), 500
    
    try:
        servin_client.restart_container(container_id)
        return jsonify({'success': True, 'message': f'Container {container_id} restarted'})
    except ServinError as e:
        return jsonify({'error': str(e)}), 500

@app.route('/api/containers/<container_id>/remove', methods=['DELETE'])
def remove_container(container_id):
    """Remove a container"""
    if not servin_client:
        return jsonify({'error': 'Servin runtime not available'}), 500
    
    try:
        servin_client.remove_container(container_id, force=True)
        return jsonify({'success': True, 'message': f'Container {container_id} removed'})
    except ServinError as e:
        return jsonify({'error': str(e)}), 500

@app.route('/api/containers/<container_id>/details', methods=['GET'])
def get_container_details(container_id):
    """Get detailed information about a container"""
    if not servin_client:
        return jsonify({'error': 'Servin runtime not available'}), 500
    
    try:
        details = servin_client.inspect_container(container_id)
        return jsonify(details)
    except ServinError as e:
        return jsonify({'error': str(e)}), 500

@app.route('/api/containers/<container_id>/logs', methods=['GET'])
def get_container_logs(container_id):
    """Get container logs"""
    if not servin_client:
        return jsonify({'error': 'Servin runtime not available'}), 500
    
    try:
        follow = request.args.get('follow', 'false').lower() == 'true'
        tail = request.args.get('tail', '100')
        logs = servin_client.get_logs(container_id, follow=follow, tail=int(tail))
        return jsonify({'logs': logs})
    except ServinError as e:
        return jsonify({'error': str(e)}), 500

@app.route('/api/containers/<container_id>/files', methods=['GET'])
def get_container_files(container_id):
    """Get container filesystem listing"""
    if not servin_client:
        return jsonify({'error': 'Servin runtime not available'}), 500
    
    try:
        path = request.args.get('path', '/')
        files = servin_client.list_files(container_id, path)
        return jsonify({'path': path, 'files': files})
    except ServinError as e:
        return jsonify({'error': str(e)}), 500

@app.route('/api/containers/<container_id>/exec', methods=['POST'])
def exec_container_command(container_id):
    """Execute command in container"""
    if not servin_client:
        return jsonify({'error': 'Servin runtime not available'}), 500
    
    try:
        data = request.get_json()
        command = data.get('command', 'sh')
        interactive = data.get('interactive', True)
        
        if interactive:
            # For interactive sessions, we'll simulate the response
            # In a real implementation, this would use WebSockets
            result = servin_client.exec_command(container_id, command)
            return jsonify({'output': result, 'exit_code': 0})
        else:
            result = servin_client.exec_command(container_id, command)
            return jsonify({'output': result, 'exit_code': 0})
    except ServinError as e:
        return jsonify({'error': str(e)}), 500

@app.route('/api/containers/<container_id>/env', methods=['GET'])
def get_container_environment(container_id):
    """Get container environment variables"""
    if not servin_client:
        return jsonify({'error': 'Servin runtime not available'}), 500
    
    try:
        env_vars = servin_client.get_environment(container_id)
        return jsonify({'environment': env_vars})
    except ServinError as e:
        return jsonify({'error': str(e)}), 500

# Image Management APIs
@app.route('/api/images', methods=['GET'])
def get_images():
    """Get list of all images"""
    if not servin_client:
        return jsonify({'error': 'Servin runtime not available'}), 500
    
    try:
        images = servin_client.list_images()
        return jsonify(images)
    except ServinError as e:
        return jsonify({'error': str(e)}), 500

@app.route('/api/images/pull', methods=['POST'])
def pull_image():
    """Pull an image from registry (import for servin)"""
    if not servin_client:
        return jsonify({'error': 'Servin runtime not available'}), 500
    
    data = request.get_json()
    if not data or 'image' not in data:
        return jsonify({'error': 'Image name required'}), 400
    
    try:
        # For servin, we'll explain that pull is not supported
        # Instead, users need to import images from tarballs
        return jsonify({'error': 'Servin does not support pulling images from registries. Use image import with tarball files instead.'}), 400
    except ServinError as e:
        return jsonify({'error': str(e)}), 500

@app.route('/api/images/import', methods=['POST'])
def import_image():
    """Import an image from tarball"""
    if not servin_client:
        return jsonify({'error': 'Servin runtime not available'}), 500
    
    data = request.get_json()
    if not data or 'tarball' not in data or 'name' not in data:
        return jsonify({'error': 'Tarball path and image name required'}), 400
    
    try:
        tarball_path = data['tarball']
        image_name = data['name']
        servin_client.import_image(tarball_path, image_name)
        return jsonify({'success': True, 'message': f'Image {image_name} imported successfully'})
    except ServinError as e:
        return jsonify({'error': str(e)}), 500

@app.route('/api/images/<image_id>/remove', methods=['DELETE'])
def remove_image(image_id):
    """Remove an image"""
    if not servin_client:
        return jsonify({'error': 'Servin runtime not available'}), 500
    
    try:
        servin_client.remove_image(image_id, force=True)
        return jsonify({'success': True, 'message': f'Image {image_id} removed'})
    except ServinError as e:
        return jsonify({'error': str(e)}), 500

# Volume Management APIs
@app.route('/api/volumes', methods=['GET'])
def get_volumes():
    """Get list of all volumes"""
    if not servin_client:
        return jsonify({'error': 'Servin runtime not available'}), 500
    
    try:
        volumes = servin_client.list_volumes()
        return jsonify(volumes)
    except ServinError as e:
        return jsonify({'error': str(e)}), 500

@app.route('/api/volumes/<volume_name>/remove', methods=['DELETE'])
def remove_volume(volume_name):
    """Remove a volume"""
    if not servin_client:
        return jsonify({'error': 'Servin runtime not available'}), 500
    
    try:
        servin_client.remove_volume(volume_name)
        return jsonify({'success': True, 'message': f'Volume {volume_name} removed'})
    except ServinError as e:
        return jsonify({'error': str(e)}), 500

@app.route('/api/volumes/create', methods=['POST'])
def create_volume():
    """Create a new volume"""
    if not servin_client:
        return jsonify({'error': 'Servin runtime not available'}), 500
    
    data = request.get_json()
    if not data or 'name' not in data:
        return jsonify({'error': 'Volume name required'}), 400
    
    try:
        volume_name = data['name']
        servin_client.create_volume(volume_name)
        return jsonify({'success': True, 'message': f'Volume {volume_name} created successfully'})
    except ServinError as e:
        return jsonify({'error': str(e)}), 500

# System Information APIs
@app.route('/api/system/info', methods=['GET'])
def get_system_info():
    """Get Servin system information"""
    if not servin_client:
        return jsonify({'error': 'Servin runtime not available'}), 500
    
    try:
        info = servin_client.info()
        return jsonify(info)
    except ServinError as e:
        return jsonify({'error': str(e)}), 500

def run_flask_app():
    """Run the Flask application"""
    app.run(host='127.0.0.1', port=5555, debug=False, use_reloader=False)

if __name__ == '__main__':
    run_flask_app()
