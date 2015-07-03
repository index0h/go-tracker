if [ -f /etc/bashrc ]; then
	. /etc/bashrc
fi

source $HOME/.bash_aliases
source $HOME/.bash_git
source $HOME/.gvm/scripts/gvm

export PATH
