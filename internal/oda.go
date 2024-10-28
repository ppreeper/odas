package internal

import (
	"bufio"
	"embed"
	"os"
	"path/filepath"
	"regexp"
)

type OdooDatabase struct {
	Name    string
	Version string
	Image   string
}

type OdooConfig struct {
	Name             string   `json:"name"`
	Version          string   `json:"version"`
	Image            string   `json:"image"`
	InstanceName     string   `json:"instance_name"`
	Repos            []string `json:"repos"`
	BaselinePackages []string `json:"baseline_packages"`
	Odoobase         []string `json:"odoobase"`
}

type QueryDef struct {
	Model    string
	Filter   string
	Offset   int
	Limit    int
	Fields   string
	Count    bool
	Username string
	Password string
}

type OdooConf struct {
	DbHost     string
	DbPort     string
	DbName     string
	DbUser     string
	DbPassword string
	DbTemplate string
	AddonsPath string
	DataDir    string
}

type ODA struct {
	Name         string
	Usage        string
	Version      string
	EmbedFS      embed.FS
	Q            QueryDef
	OdooRepos    []string
	OdooVersions []string
	OdooDatabase OdooDatabase
	OdooConfigs  []OdooConfig
	OdooConf     OdooConf
}

func NewODA(name, usage, version string, embedFS embed.FS) *ODA {
	return &ODA{
		Name:         name,
		Usage:        usage,
		Version:      version,
		EmbedFS:      embedFS,
		Q:            QueryDef{},
		OdooRepos:    []string{"odoo", "enterprise", "design-themes", "industry"},
		OdooVersions: []string{"18.0", "17.0", "16.0", "15.0"},
		OdooDatabase: OdooDatabase{Name: "db", Version: "17", Image: "debian/12"},
		OdooConfigs: []OdooConfig{
			{
				Name:         "15",
				Version:      "15.0",
				Image:        "ubuntu/22.04",
				InstanceName: "odoo-15-0",
				Repos:        []string{"odoo", "enterprise", "design-themes"},
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
				Repos:        []string{"odoo", "enterprise", "design-themes", "industry"},
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
				Repos:        []string{"odoo", "enterprise", "design-themes", "industry"},
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
				Repos:        []string{"odoo", "enterprise", "design-themes", "industry"},
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
		},
	}
}

func ReadConfValue(conffile, key, def string) string {
	c, err := os.Open(conffile)
	if err != nil {
		return def
	}
	defer func() {
		if err := c.Close(); err != nil {
			panic(err)
		}
	}()
	scanner := bufio.NewScanner(c)
	for scanner.Scan() {
		line := scanner.Text()
		re := regexp.MustCompile(`^` + key + ` = (.+)$`)
		if re.MatchString(line) {
			match := re.FindStringSubmatch(line)
			return match[1]
		}
	}
	return def
}

func (o *ODA) GetOdooConf() *ODA {
	odooconf := filepath.Join("/", "opt", "odoo", "conf", "odoo.conf")
	o.OdooConf.DbHost = ReadConfValue(odooconf, "db_host", "localhost")
	o.OdooConf.DbPort = ReadConfValue(odooconf, "db_port", "5432")
	o.OdooConf.DbName = ReadConfValue(odooconf, "db_name", "odoo")
	o.OdooConf.DbUser = ReadConfValue(odooconf, "db_user", "odoo")
	o.OdooConf.DbPassword = ReadConfValue(odooconf, "db_password", "odoo")
	o.OdooConf.DbTemplate = ReadConfValue(odooconf, "db_template", "template0")
	o.OdooConf.AddonsPath = ReadConfValue(odooconf, "addons_path", "/opt/odoo/addons")
	o.OdooConf.DataDir = ReadConfValue(odooconf, "data_dir", "/opt/odoo/data")
	return o
}
