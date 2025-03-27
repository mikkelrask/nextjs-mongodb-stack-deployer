# nextjs-mongodb-stack-deployer
Deployer made in Go for my [nextjs-mongodb-stack](/mikkelrask/nextjs-mongodb-stack) thing, that, as the name implies, builds docker image for a nextjs website with a mongodb database and spins up the stack purely from environment variables.

## Download
Go to [Releases](https://github.com/mikkelrask/nextjs-mongodb-stack-deployer/releases) and download the latest version for your OS/Arch.  
Make it executeable and run it with `./deployer-<OS>-<architecture>` 

Can be placed in your `$PATH` for system wide callability.

## Build yourself
```
git clone https://github.com/mikkelrask/nextjs-mongodb-stack-deployer.git deployer
cd deployer
go mod tidy
go build -o deployer
```

## Usage
Invoke the deployer and fill out the information needed to start your project.
```
./deployer
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
