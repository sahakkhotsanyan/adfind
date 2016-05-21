# adfind
Admin Panel Finder
##Dependen's php
##Installing php and git on arch
sudo pacman -S php
sudo pacman -S git
##Installing php and git on debian
sudo apt-get install php7.0-cli
sudo apt-get install git
##installing adfind
sudo git clone https://github.com/sahakkhotsanyan/adfind.git
cd adfind*
sudo cp adfind /bin/adfind
sudo chmod +x /bin/adfind
##Usage
adfind http://example.com
