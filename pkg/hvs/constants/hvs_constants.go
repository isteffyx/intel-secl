package constants

import "time"

// 2 things done on this file:
//   1. remove unused constants that does not have clear meaning
//   2. categorize stuffs

// general HVS constants
const (
	ServiceName     = "HVS"
	ServiceDir      = "hvs/"
	OldServiceName  = "mtwilson"
	ApiVersion      = "/v2"
	ServiceUserName = "hvs"

	// Timestamp operations
	HostStatusDateFormat = "2006-01-02 15:04:05"

	// service remove command
	ServiceRemoveCmd = "systemctl disable hvs"

	HVSParamDateFormat = "2006-01-02 15:04:05"
)

// these are used only when uninstalling service
const (
	HomeDir      = "/opt/" + ServiceDir
	RunDirPath   = "/run/" + ServiceDir
	ExecLinkPath = "/usr/bin/" + ServiceUserName
	LogDir       = "/var/log/" + ServiceDir
)

// file and directory constants
const (
	DefaultConfigFilePath = "/etc/hvs/config.yaml"
	ConfigDir             = "/etc/hvs/"
	ConfigFile            = "config"

	// certificates' path
	TrustedJWTSigningCertsDir = ConfigDir + "certs/trustedjwt/"
	TrustedCaCertsDir         = ConfigDir + "certs/trustedca/"
	TrustedRootCACertsDir     = TrustedCaCertsDir + "root/"

	TrustedKeysDir = ConfigDir + "trusted-keys/"

	// saml key and cert
	SAMLCertFile = TrustedCaCertsDir + "saml-cert.pem"
	SAMLKeyFile  = TrustedKeysDir + "saml.key"

	// flavor signing key and cert
	FlavorSigningCertFile = TrustedCaCertsDir + "flavor-signing.pem"
	FlavorSigningKeyFile  = TrustedKeysDir + "flavor-signing.key"

	// privacy ca key and cert
	PrivacyCACertFile = TrustedCaCertsDir + "privacy-ca/privacy-ca-cert.pem"
	PrivacyCAKeyFile  = TrustedKeysDir + "privacy-ca.key"

	//TODO remove or dont use temporary files
	AikRequestsDir = HomeDir + "privacyca-aik-requests/"
	//TODO use EndorsementCA file after implementation of create_endorsement_ca setup task
	EndorsementCACertDir      = ConfigDir + "certs/endorsement/"
	EndorsementCACertFile     = EndorsementCACertDir + "EndorsementCA-external.pem" //External ECA
	SelfEndorsementCACertFile = EndorsementCACertDir + "EndorsementCA.pem" //Self signed ECA
	EndorsementCAKeyFile      = TrustedKeysDir + "endorsement-ca.key"

	TagCACertFile = TrustedCaCertsDir + "tag-ca-cert.pem"
	TagCAKeyFile  = TrustedKeysDir + "tag-ca.key"

	// default locations for tls certificate and key
	DefaultTLSKeyFile  = ConfigDir + "tls.key"
	DefaultTLSCertFile = ConfigDir + "tls-cert.pem"
)

// tls constants
const (
	DefaultHvsTlsCn           = "HVS TLS Certificate"
	DefaultHvsTlsSan          = "127.0.0.1,localhost"
	DefaultKeyAlgorithm       = "rsa"
	DefaultKeyAlgorithmLength = 3072
)

// aik and privacy ca constants
const (
	AIKCertValidity                = 1
	DefaultPrivacyCACertValidity   = 5
	HostSigningKeyCertificateCN    = "Signing_Key_Certificate"
	HostBindingKeyCertificateCN    = "Binding_Key_Certificate"
	DefaultPrivacyCaIdentityIssuer = "hvs-pca-aik"
)

// general constants for certificates
const (
	DefaultCN           = "HVS Certificate"
	DefaultSANList      = "127.0.0.1,localhost"
	DefaultCertIssuer   = "intel-secl"
	DefaultCertValidity = 5

	NAString = "N/A"

	DefaultTagCertValiditySeconds = 60 * 60 * 24 * 365
)

// server costants
const (
	DefaultHVSListenerPort   = 8443
	DefaultReadTimeout       = 30 * time.Second
	DefaultReadHeaderTimeout = 10 * time.Second
	DefaultWriteTimeout      = 10 * time.Second
	DefaultIdleTimeout       = 10 * time.Second
	DefaultMaxHeaderBytes    = 1 << 20
)

// db constants
const (
	DBTypePostgres = "postgres"

	DefaultDbConnRetryAttempts  = 4
	DefaultDbConnRetryTime      = 1
	DefaultSearchResultRowLimit = 10000

	//Postgres connection SslModes
	SslModeAllow      = "allow"
	SslModePrefer     = "prefer"
	SslModeVerifyCa   = "verify-ca"
	SslModeRequire    = "require"
	SslModeVerifyFull = "verify-full"
)

// log constants
const (
	DefaultLogEntryMaxlength = 300
)

// jwt constants
const (
	JWTCertsCacheTime = "1m"
)
