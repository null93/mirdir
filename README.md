# mirdir (mirror directory)

## About

...

## Install

<details>
  <summary>Darwin</summary>

  ### Intel & ARM

  ```shell
  brew tap null93/tap
  brew install mirdir
  ```
</details>

<details>
  <summary>Debian</summary>

  ### amd64

  ```shell
  curl -sL -o ./mirdir_0.0.1_amd64.deb https://github.com/null93/mirdir/releases/download/0.0.1/mirdir_0.0.1_amd64.deb
  sudo dpkg -i ./mirdir_0.0.1_amd64.deb
  rm ./mirdir_0.0.1_amd64.deb
  ```

  ### arm64

  ```shell
  curl -sL -o ./mirdir_0.0.1_arm64.deb https://github.com/null93/mirdir/releases/download/0.0.1/mirdir_0.0.1_arm64.deb
  sudo dpkg -i ./mirdir_0.0.1_arm64.deb
  rm ./mirdir_0.0.1_arm64.deb
  ```
</details>

<details>
  <summary>Red Hat</summary>

  ### aarch64

  ```shell
  rpm -i https://github.com/null93/mirdir/releases/download/0.0.1/mirdir-0.0.1-1.aarch64.rpm
  ```

  ### x86_64

  ```shell
  rpm -i https://github.com/null93/mirdir/releases/download/0.0.1/mirdir-0.0.1-1.x86_64.rpm
  ```
</details>

<details>
  <summary>Alpine</summary>

  ### aarch64

  ```shell
  curl -sL -o ./mirdir_0.0.1_aarch64.apk https://github.com/null93/mirdir/releases/download/0.0.1/mirdir_0.0.1_aarch64.apk
  apk add --allow-untrusted ./mirdir_0.0.1_aarch64.apk
  rm ./mirdir_0.0.1_aarch64.apk
  ```

  ### x86_64

  ```shell
  curl -sL -o ./mirdir_0.0.1_x86_64.apk https://github.com/null93/mirdir/releases/download/0.0.1/mirdir_0.0.1_x86_64.apk
  apk add --allow-untrusted ./mirdir_0.0.1_x86_64.apk
  rm ./mirdir_0.0.1_x86_64.apk
  ```
</details>

<details>
  <summary>Arch</summary>

  ### aarch64

  ```shell
  curl -sL -o ./mirdir-0.0.1-1-aarch64.pkg.tar.zst https://github.com/null93/mirdir/releases/download/0.0.1/mirdir-0.0.1-1-aarch64.pkg.tar.zst
  sudo pacman -U ./mirdir-0.0.1-1-aarch64.pkg.tar.zst
  rm ./mirdir-0.0.1-1-aarch64.pkg.tar.zst
  ```

  ### x86_64

  ```shell
  curl -sL -o ./mirdir-0.0.1-1-x86_64.pkg.tar.zst https://github.com/null93/mirdir/releases/download/0.0.1/mirdir-0.0.1-1-x86_64.pkg.tar.zst
  sudo pacman -U ./mirdir-0.0.1-1-x86_64.pkg.tar.zst
  rm ./mirdir-0.0.1-1-x86_64.pkg.tar.zst
  ```
</details>

## Usage

```shell
$ tree /opt/confs/nginx

/opt/confs/nginx
├── conf.d
│   ├── bar.conf
│   ├── baz.conf
│   └── foo.conf
├── nginx.conf
├── sites-available
│   └── custom-site-[USER].conf.tpl
└── sites-enabled
    └── custom-site-[USER].conf -> ../sites-available/custom-site-[USER].conf

4 directories, 6 files
```

```shell
$ mirdir /opt/confs/nginx/ /etc/nginx/
```

```shell
$ tree /etc/nginx

/etc/nginx
├── conf.d
│   ├── bar.conf
│   ├── baz.conf
│   └── foo.conf
├── nginx.conf
├── sites-available
│   └── custom-site-raffi.conf
└── sites-enabled
    └── custom-site-raffi.conf -> ../sites-available/custom-site-raffi.conf

4 directories, 6 files
```

```shell
$ mirdir /opt/confs/nginx/ /etc/nginx/ --dry-run

-rwxrwxrwx  501:20   /etc/nginx/
-rwxrwxrwx  501:20   /etc/nginx/conf.d/
-rwxrwxrwx  501:20   /etc/nginx/sites-available/
-rwxrwxrwx  501:20   /etc/nginx/sites-enabled/
-rw-rw-rw-  501:20   /etc/nginx/conf.d/bar.conf
add_header X-Bar "bar";

-rw-rw-rw-  501:20   /etc/nginx/conf.d/baz.conf
add_header X-Baz "baz";

-rw-rw-rw-  501:20   /etc/nginx/conf.d/foo.conf
add_header X-Foo "foo";

-rw-rw-rw-  501:20   /etc/nginx/nginx.conf
user nginx;
worker_processes auto;

error_log /var/log/nginx/error.log notice;
pid /var/run/nginx.pid;

events {
	worker_connections 1024;
}

http {
	default_type application/octet-stream;
	include ./sites-enabled/*.conf;

	server {
		listen 80;
		server_name _;
		root /var/www/html;
	}

}

-rw-rw-rw-  501:20   /etc/nginx/sites-available/custom-site-raffi.conf
server {
	listen 80;
	server_name raffi.local;
	root /home/raffi/public;
	index index.html;
	include ../conf.d/*.conf;
}

-rw-rw-rw-  501:20   /etc/nginx/sites-enabled/custom-site-raffi.conf -> ../sites-available/custom-site-raffi.conf
```

