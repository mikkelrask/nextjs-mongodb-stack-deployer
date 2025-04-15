# nextjs-mongodb-stack-deployer
Deployer made in Go for my [nextjs-mongodb-stack](/mikkelrask/nextjs-mongodb-stack) thing, that, as the name implies, builds docker image for a nextjs website with a mongodb database and spins up the stack purely from environment variables.

## Download
Go to [Releases](https://github.com/mikkelrask/nextjs-mongodb-stack-deployer/releases) and download the latest version for your OS/Arch.  
Make it executeable, and you can call it with `./deployer-<OS>-<architecture>` 


### Download example
This would how you get the `v.0.1.2` version, to your `~/.local/bin/` directory and make it executeable systemwide by your user. 
```bash
wget https://github.com/mikkelrask/nextjs-mongodb-stack-deployer/releases/download/v0.1.2/deployer-<YOUR-OS>-<YOUR-ARCHITECTURE> -O ~/.local/bin/deployer
chmod +x ~/.local/bin/deployer
deployer
🚀 DPLOY - BUILD YOUR WEBSITE

Fill the needed data to deploy your webapp - Enter keeps the suggested default

? NextJS webapp repository URL: https://github.com/SiddharthaMaity/nextjs-15-starter-core.git
? Is the repository private? No
? Do you want to specify a port for NextJS?  3001
? Do you want to import a database dump? Yes
? Database dump path:  ~/Downloads/production/production
📁 Copying the database dump to the dump directory...
? MongoDB Database name? production
? MongoDB Username:  admin
....
....
```


## Build with Go
If you prefer to build it manually or make changes to it before using, you can do so like any other Go projects:
```
git clone https://github.com/mikkelrask/nextjs-mongodb-stack-deployer.git deployer
cd deployer
go mod tidy
go build -o deployer
```
After building `deployer` will be in the root of the repo. 

Can be placed anywhere in your `$PATH` for system wide callability.

## Usage
Invoke the deployer and fill out the information needed to start your project.
```
deployer
🚀 DPLOY - BUILD YOUR WEBSITE

Fill the needed data to deploy your webapp - Enter keeps the suggested default

? NextJS webapp repository URL: https://github.com/SiddharthaMaity/nextjs-15-starter-core.git
? Is the repository private? No
? Do you want to specify a port for NextJS?  3001
? Do you want to import a database dump? Yes
? Database dump path:  ~/Downloads/production/production
📁 Copying the database dump to the dump directory...
? MongoDB Database name? production
? MongoDB Username:  admin
....
....
```
