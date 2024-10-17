package internal

import "embed"

//go:embed templates/*
var embedFS embed.FS

var OdooDatabase = struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Image   string `json:"image"`
}{
	Name:    "db",
	Version: "17",
	Image:   "debian/12",
}

var OdooRepos = []string{"odoo", "enterprise", "design-themes", "industry"}

var OdooVersions = []string{"18.0", "17.0", "16.0", "15.0"}

type OdooConfig struct {
	Name             string   `json:"name"`
	Version          string   `json:"version"`
	Image            string   `json:"image"`
	InstanceName     string   `json:"instance_name"`
	BaselinePackages []string `json:"baseline_packages"`
	Odoobase         []string `json:"odoobase"`
}

var OdooConfigs = []OdooConfig{
	{
		Name:         "15",
		Version:      "15.0",
		Image:        "ubuntu/22.04",
		InstanceName: "odoo-15-0",
		BaselinePackages: []string{
			"apt-transport-https", "apt-utils", "bzip2", "ca-certificates", "curl",
			"dirmngr", "git", "gnupg", "inetutils-ping", "libgnutls-dane0", "libgts-bin",
			"libpaper-utils", "locales", "lsb-release", "nodejs", "npm", "odaserver",
			"openssh-server", "postgresql-common", "python3", "python3-full",
			"shared-mime-info", "sudo", "unzip", "vim", "wget", "xz-utils", "zip", "zstd",
		},
		Odoobase: []string{
			"fonts-liberation", "fonts-noto", "fonts-noto-cjk", "fonts-noto-mono",
			"geoip-database", "gsfonts", "python3-babel", "python3-chardet",
			"python3-cryptography", "python3-cups", "python3-dateutil",
			"python3-decorator", "python3-docutils", "python3-feedparser",
			"python3-freezegun", "python3-geoip2", "python3-gevent", "python3-googleapi",
			"python3-greenlet", "python3-html2text", "python3-idna", "python3-jinja2",
			"python3-ldap", "python3-libsass", "python3-lxml", "python3-markupsafe",
			"python3-num2words", "python3-odf", "python3-ofxparse", "python3-olefile",
			"python3-openssl", "python3-paramiko", "python3-passlib", "python3-pdfminer",
			"python3-phonenumbers", "python3-pil", "python3-pip", "python3-polib",
			"python3-psutil", "python3-psycopg2", "python3-pydot", "python3-pylibdmtx",
			"python3-pyparsing", "python3-pypdf2", "python3-qrcode", "python3-renderpm",
			"python3-reportlab", "python3-reportlab-accel", "python3-requests",
			"python3-rjsmin", "python3-serial", "python3-setuptools", "python3-stdnum",
			"python3-tz", "python3-urllib3", "python3-usb", "python3-vobject",
			"python3-werkzeug", "python3-xlrd", "python3-xlsxwriter", "python3-xlwt",
			"python3-zeep",
		},
	},
	{
		Name:         "16",
		Version:      "16.0",
		Image:        "ubuntu/22.04",
		InstanceName: "odoo-16-0",
		BaselinePackages: []string{
			"apt-transport-https", "apt-utils", "bzip2", "ca-certificates", "curl",
			"dirmngr", "git", "gnupg", "inetutils-ping", "libgnutls-dane0", "libgts-bin",
			"libpaper-utils", "locales", "lsb-release", "nodejs", "npm", "odaserver",
			"openssh-server", "postgresql-common", "python3", "python3-full",
			"shared-mime-info", "sudo", "unzip", "vim", "wget", "xz-utils", "zip", "zstd",
		},
		Odoobase: []string{
			"fonts-liberation", "fonts-noto", "fonts-noto-cjk", "fonts-noto-mono",
			"geoip-database", "gsfonts", "python3-babel", "python3-chardet",
			"python3-cryptography", "python3-cups", "python3-dateutil",
			"python3-decorator", "python3-docutils", "python3-feedparser",
			"python3-freezegun", "python3-geoip2", "python3-gevent", "python3-googleapi",
			"python3-greenlet", "python3-html2text", "python3-idna", "python3-jinja2",
			"python3-ldap", "python3-libsass", "python3-lxml", "python3-markupsafe",
			"python3-num2words", "python3-odf", "python3-ofxparse", "python3-olefile",
			"python3-openssl", "python3-paramiko", "python3-passlib", "python3-pdfminer",
			"python3-phonenumbers", "python3-pil", "python3-pip", "python3-polib",
			"python3-psutil", "python3-psycopg2", "python3-pydot", "python3-pylibdmtx",
			"python3-pyparsing", "python3-pypdf2", "python3-qrcode", "python3-renderpm",
			"python3-reportlab", "python3-reportlab-accel", "python3-requests",
			"python3-rjsmin", "python3-serial", "python3-setuptools", "python3-stdnum",
			"python3-tz", "python3-urllib3", "python3-usb", "python3-vobject",
			"python3-werkzeug", "python3-xlrd", "python3-xlsxwriter", "python3-xlwt",
			"python3-zeep",
		},
	},
	{
		Name:         "17",
		Version:      "17.0",
		Image:        "ubuntu/22.04",
		InstanceName: "odoo-17-0",
		BaselinePackages: []string{
			"apt-transport-https", "apt-utils", "bzip2", "ca-certificates", "curl",
			"dirmngr", "git", "gnupg", "inetutils-ping", "libgnutls-dane0", "libgts-bin",
			"libpaper-utils", "locales", "lsb-release", "nodejs", "npm", "odaserver",
			"openssh-server", "postgresql-common", "python3", "python3-full",
			"shared-mime-info", "sudo", "unzip", "vim", "wget", "xz-utils", "zip", "zstd",
		},
		Odoobase: []string{
			"fonts-liberation", "fonts-noto", "fonts-noto-cjk", "fonts-noto-mono",
			"geoip-database", "gsfonts", "python3-babel", "python3-chardet",
			"python3-cryptography", "python3-cups", "python3-dateutil",
			"python3-decorator", "python3-docutils", "python3-feedparser",
			"python3-freezegun", "python3-geoip2", "python3-gevent", "python3-googleapi",
			"python3-greenlet", "python3-html2text", "python3-idna", "python3-jinja2",
			"python3-ldap", "python3-libsass", "python3-lxml", "python3-markupsafe",
			"python3-num2words", "python3-odf", "python3-ofxparse", "python3-olefile",
			"python3-openssl", "python3-paramiko", "python3-passlib", "python3-pdfminer",
			"python3-phonenumbers", "python3-pil", "python3-pip", "python3-polib",
			"python3-psutil", "python3-psycopg2", "python3-pydot", "python3-pylibdmtx",
			"python3-pyparsing", "python3-pypdf2", "python3-qrcode", "python3-renderpm",
			"python3-reportlab", "python3-reportlab-accel", "python3-requests",
			"python3-rjsmin", "python3-serial", "python3-setuptools", "python3-stdnum",
			"python3-tz", "python3-urllib3", "python3-usb", "python3-vobject",
			"python3-werkzeug", "python3-xlrd", "python3-xlsxwriter", "python3-xlwt",
			"python3-zeep",
		},
	},
	{
		Name:         "18",
		Version:      "18.0",
		Image:        "ubuntu/24.04",
		InstanceName: "odoo-18-0",
		BaselinePackages: []string{
			"apt-transport-https", "apt-utils", "bzip2", "ca-certificates", "curl",
			"dirmngr", "git", "gnupg", "inetutils-ping", "libgnutls-dane0", "libgts-bin",
			"libpaper-utils", "locales", "lsb-release", "nodejs", "npm", "odaserver",
			"openssh-server", "postgresql-common", "python3", "python3-full",
			"shared-mime-info", "sudo", "unzip", "vim", "wget", "xz-utils", "zip", "zstd",
		},
		Odoobase: []string{
			"fonts-liberation", "fonts-noto", "fonts-noto-cjk", "fonts-noto-mono",
			"geoip-database", "gsfonts", "python3-asn1crypto", "python3-babel",
			"python3-cbor2", "python3-chardet", "python3-cryptography", "python3-cups",
			"python3-dateutil", "python3-decorator", "python3-docutils",
			"python3-feedparser", "python3-freezegun", "python3-geoip2",
			"python3-gevent", "python3-googleapi", "python3-greenlet",
			"python3-html2text", "python3-idna", "python3-jinja2", "python3-ldap",
			"python3-libsass", "python3-lxml", "python3-lxml-html-clean",
			"python3-markupsafe", "python3-num2words", "python3-odf", "python3-ofxparse",
			"python3-olefile", "python3-openpyxl", "python3-openssl", "python3-paramiko",
			"python3-passlib", "python3-pdfminer", "python3-phonenumbers", "python3-pil",
			"python3-pip", "python3-polib", "python3-psutil", "python3-psycopg2",
			"python3-pydot", "python3-pylibdmtx", "python3-pyparsing", "python3-pypdf2",
			"python3-qrcode", "python3-renderpm", "python3-reportlab",
			"python3-rl-renderpm", "python3-reportlab-accel", "python3-requests",
			"python3-rjsmin", "python3-serial", "python3-setuptools", "python3-stdnum",
			"python3-tz", "python3-urllib3", "python3-usb", "python3-vobject",
			"python3-werkzeug", "python3-xlrd", "python3-xlsxwriter", "python3-xlwt",
			"python3-zeep",
		},
	},
}
