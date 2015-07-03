Vagrant.require_version ">= 1.5"

Vagrant.configure("2") do |config|

    config.vm.provider :virtualbox do |v|
        v.name = "go-tracker"
        v.customize [
            "modifyvm", :id,
            "--name", "go-tracker",
            "--memory", 1024,
            "--natdnshostresolver1", "on",
            "--cpus", 1,
        ]
    end

    config.vm.box = "ubuntu/trusty64"

    config.vm.network :private_network, ip: "192.168.100.126"
    config.ssh.forward_agent = true

    config.vm.provision "ansible" do |ansible|
        ansible.playbook = "ansible/playbook.yml"
        ansible.inventory_path = "ansible/inventories/dev"
        ansible.limit = 'all'
    end

    config.vm.synced_folder "./", "/vagrant"
    config.vm.synced_folder "./", "/home/vagrant/work"
end
