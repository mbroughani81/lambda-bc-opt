#!/bin/bash

# Function to install Docker on Ubuntu
install_docker_ubuntu() {
    echo "Updating package information..."
    sudo apt-get update -y

    echo "Installing required dependencies..."
    sudo apt-get install apt-transport-https ca-certificates curl software-properties-common -y

    echo "Adding Docker’s official GPG key..."
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg

    echo "Adding Docker repository..."
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

    echo "Updating package information with Docker packages..."
    sudo apt-get update -y

    echo "Installing Docker..."
    sudo apt-get install docker-ce docker-ce-cli containerd.io -y

    echo "Starting and enabling Docker..."
    sudo systemctl start docker
    sudo systemctl enable docker

    echo "Docker installed successfully on Ubuntu!"
}

# Function to install Docker Compose
install_docker_compose() {
    echo "Installing Docker Compose..."

    # Download the latest version of Docker Compose (replace version if needed)
    DOCKER_COMPOSE_VERSION="2.20.2"
    sudo curl -L "https://github.com/docker/compose/releases/download/v${DOCKER_COMPOSE_VERSION}/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose

    # Make the Docker Compose binary executable
    sudo chmod +x /usr/local/bin/docker-compose

    # Verify installation
    docker-compose --version

    echo "Docker Compose installed successfully!"
}

# Install Docker and Docker Compose
install_docker_ubuntu
install_docker_compose

# Add the current user to the docker group to avoid using sudo for Docker commands
echo "Adding user $(whoami) to the docker group..."
sudo usermod -aG docker $(whoami)

echo "Docker and Docker Compose installation complete. You may need to log out and back in for group changes to take effect."
