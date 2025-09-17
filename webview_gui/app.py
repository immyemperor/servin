"""
Servin Runtime
Flask API server for managing Servin containers, images, and volumes
"""

import os
import sys
import threading
import time
import subprocess
import json
from flask import Flask, jsonify, request, render_template, send_from_directory
from flask_cors import CORS
from flask_socketio import SocketIO, emit, disconnect
from servin_client import ServinClient, ServinError

app = Flask(__name__)
app.config['SECRET_KEY'] = 'servin-gui-secret-key'
CORS(app)  # Enable CORS for all routes
socketio = SocketIO(app, cors_allowed_origins="*")

# Store active log streaming processes
active_log_streams = {}
active_exec_sessions = {}

# Initialize Servin client
try:
    from servin_client import ServinClient, ServinError
    servin_client = ServinClient()
    # Test connection
    servin_client.ping()
    print("Successfully connected to Servin runtime")
except (ServinError, Exception) as e:
    print(f"Failed to connect to Servin runtime: {e}")
    print("Please ensure the servin binary is available and working properly")
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
        return jsonify(files)
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

# WebSocket Event Handlers for Real-time Features

@socketio.on('connect')
def handle_connect():
    """Handle client connection"""
    print(f"Client connected: {request.sid}")
    emit('status', {'message': 'Connected to Servin GUI'})

@socketio.on('disconnect')
def handle_disconnect():
    """Handle client disconnection"""
    print(f"Client disconnected: {request.sid}")
    # Clean up any active streams for this client
    cleanup_client_streams(request.sid)

@socketio.on('start_logs')
def handle_start_logs(data):
    """Start streaming logs for a container"""
    container_id = data.get('container_id')
    if not container_id:
        emit('error', {'message': 'Container ID required'})
        return
    
    if not servin_client:
        emit('error', {'message': 'Servin runtime not available'})
        return
    
    try:
        # Stop any existing log stream for this container and client
        stream_key = f"{request.sid}:{container_id}"
        if stream_key in active_log_streams:
            active_log_streams[stream_key]['stop'] = True
        
        # Start new log streaming thread
        active_log_streams[stream_key] = {'stop': False}
        thread = threading.Thread(
            target=stream_logs_thread,
            args=(container_id, request.sid, stream_key)
        )
        thread.daemon = True
        thread.start()
        
        emit('logs_started', {'container_id': container_id})
    except Exception as e:
        emit('error', {'message': f'Failed to start log stream: {str(e)}'})

@socketio.on('stop_logs')
def handle_stop_logs(data):
    """Stop streaming logs for a container"""
    container_id = data.get('container_id')
    if not container_id:
        emit('error', {'message': 'Container ID required'})
        return
    
    stream_key = f"{request.sid}:{container_id}"
    if stream_key in active_log_streams:
        active_log_streams[stream_key]['stop'] = True
        del active_log_streams[stream_key]
        emit('logs_stopped', {'container_id': container_id})

@socketio.on('start_exec')
def handle_start_exec(data):
    """Start an exec session for a container"""
    container_id = data.get('container_id')
    shell = data.get('shell', '/bin/sh')
    
    if not container_id:
        emit('error', {'message': 'Container ID required'})
        return
    
    if not servin_client:
        emit('error', {'message': 'Servin runtime not available'})
        return
    
    try:
        # Stop any existing exec session for this container and client
        session_key = f"{request.sid}:{container_id}"
        if session_key in active_exec_sessions:
            active_exec_sessions[session_key]['stop'] = True
        
        # Start new exec session thread
        active_exec_sessions[session_key] = {'stop': False}
        thread = threading.Thread(
            target=exec_session_thread,
            args=(container_id, shell, request.sid, session_key)
        )
        thread.daemon = True
        thread.start()
        
        emit('exec_started', {'container_id': container_id, 'shell': shell})
    except Exception as e:
        emit('error', {'message': f'Failed to start exec session: {str(e)}'})

@socketio.on('exec_input')
def handle_exec_input(data):
    """Send input to an exec session"""
    container_id = data.get('container_id')
    command = data.get('command', '')
    
    if not container_id:
        emit('error', {'message': 'Container ID required'})
        return
    
    session_key = f"{request.sid}:{container_id}"
    if session_key not in active_exec_sessions:
        emit('error', {'message': 'No active exec session'})
        return
    
    try:
        # Add command to the session's input queue
        session = active_exec_sessions[session_key]
        if 'input_queue' not in session:
            session['input_queue'] = []
        session['input_queue'].append(command)
    except Exception as e:
        emit('error', {'message': f'Failed to send input: {str(e)}'})

@socketio.on('stop_exec')
def handle_stop_exec(data):
    """Stop an exec session"""
    container_id = data.get('container_id')
    if not container_id:
        emit('error', {'message': 'Container ID required'})
        return
    
    session_key = f"{request.sid}:{container_id}"
    if session_key in active_exec_sessions:
        active_exec_sessions[session_key]['stop'] = True
        del active_exec_sessions[session_key]
        emit('exec_stopped', {'container_id': container_id})

def stream_logs_thread(container_id, client_sid, stream_key):
    """Thread function to stream container logs"""
    try:
        # Get initial logs
        logs = servin_client.get_logs(container_id, follow=False, tail=100)
        socketio.emit('log_data', {
            'container_id': container_id,
            'data': logs,
            'type': 'initial'
        }, room=client_sid)
        
        # Start following logs
        log_process = None
        try:
            # Use servin logs command with follow
            cmd = ['servin', 'logs', '-f', container_id]
            log_process = subprocess.Popen(
                cmd,
                stdout=subprocess.PIPE,
                stderr=subprocess.STDOUT,
                universal_newlines=True,
                bufsize=1
            )
            
            while not active_log_streams.get(stream_key, {}).get('stop', True):
                if log_process.poll() is not None:
                    # Process ended
                    break
                
                line = log_process.stdout.readline()
                if line:
                    socketio.emit('log_data', {
                        'container_id': container_id,
                        'data': line.rstrip(),
                        'type': 'stream'
                    }, room=client_sid)
                else:
                    time.sleep(0.1)
                    
        except Exception as e:
            socketio.emit('error', {
                'message': f'Log streaming error: {str(e)}'
            }, room=client_sid)
        finally:
            if log_process:
                log_process.terminate()
                log_process.wait()
            
    except Exception as e:
        socketio.emit('error', {
            'message': f'Failed to stream logs: {str(e)}'
        }, room=client_sid)
    finally:
        # Clean up
        if stream_key in active_log_streams:
            del active_log_streams[stream_key]

def exec_session_thread(container_id, shell, client_sid, session_key):
    """Thread function to handle container exec session"""
    try:
        # Start exec session
        socketio.emit('exec_output', {
            'container_id': container_id,
            'data': f'Starting {shell} session in container {container_id}...\n',
            'type': 'system'
        }, room=client_sid)
        
        exec_process = None
        try:
            # Use servin exec command
            cmd = ['servin', 'exec', '-it', container_id, shell]
            exec_process = subprocess.Popen(
                cmd,
                stdin=subprocess.PIPE,
                stdout=subprocess.PIPE,
                stderr=subprocess.STDOUT,
                universal_newlines=True,
                bufsize=1
            )
            
            socketio.emit('exec_output', {
                'container_id': container_id,
                'data': f'{shell}$ ',
                'type': 'prompt'
            }, room=client_sid)
            
            while not active_exec_sessions.get(session_key, {}).get('stop', True):
                session = active_exec_sessions.get(session_key, {})
                
                # Check for input commands
                if 'input_queue' in session and session['input_queue']:
                    command = session['input_queue'].pop(0)
                    if exec_process.stdin:
                        exec_process.stdin.write(command + '\n')
                        exec_process.stdin.flush()
                    
                    # Echo the command
                    socketio.emit('exec_output', {
                        'container_id': container_id,
                        'data': command + '\n',
                        'type': 'input'
                    }, room=client_sid)
                
                # Check for output
                if exec_process.poll() is not None:
                    # Process ended
                    break
                
                # Read output (non-blocking)
                try:
                    import select
                    if select.select([exec_process.stdout], [], [], 0.1)[0]:
                        line = exec_process.stdout.readline()
                        if line:
                            socketio.emit('exec_output', {
                                'container_id': container_id,
                                'data': line,
                                'type': 'output'
                            }, room=client_sid)
                except:
                    # Fallback for platforms without select
                    time.sleep(0.1)
                    
        except Exception as e:
            socketio.emit('exec_output', {
                'container_id': container_id,
                'data': f'Exec session error: {str(e)}\n',
                'type': 'error'
            }, room=client_sid)
        finally:
            if exec_process:
                exec_process.terminate()
                exec_process.wait()
                
            socketio.emit('exec_output', {
                'container_id': container_id,
                'data': 'Exec session ended.\n',
                'type': 'system'
            }, room=client_sid)
            
    except Exception as e:
        socketio.emit('error', {
            'message': f'Failed to run exec session: {str(e)}'
        }, room=client_sid)
    finally:
        # Clean up
        if session_key in active_exec_sessions:
            del active_exec_sessions[session_key]

def cleanup_client_streams(client_sid):
    """Clean up all active streams for a disconnected client"""
    streams_to_remove = []
    sessions_to_remove = []
    
    for stream_key in active_log_streams:
        if stream_key.startswith(f"{client_sid}:"):
            active_log_streams[stream_key]['stop'] = True
            streams_to_remove.append(stream_key)
    
    for session_key in active_exec_sessions:
        if session_key.startswith(f"{client_sid}:"):
            active_exec_sessions[session_key]['stop'] = True
            sessions_to_remove.append(session_key)
    
    for stream_key in streams_to_remove:
        if stream_key in active_log_streams:
            del active_log_streams[stream_key]
            
    for session_key in sessions_to_remove:
        if session_key in active_exec_sessions:
            del active_exec_sessions[session_key]

def run_flask_app():
    """Run the Flask application with SocketIO support"""
    socketio.run(app, host='127.0.0.1', port=5555, debug=False, use_reloader=False)

if __name__ == '__main__':
    run_flask_app()
