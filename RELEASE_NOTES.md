## 1.0.0-beta.1

Initial release (public beta)

FEATURES:

  * Commands
    * Cloud Servers (`servers`)
      * `instance`
        * list/create/get/update/delete/resize/reboot/rebuild
        * create/get/update/delete metdata
        * list addresses and list addresses by network
        * boot from volume (via the `block-device` flag for `create`)
      * `flavor`
        * list/get
      * `image`
        * list/get
      * `keypair`
        * generate/upload/get/delete
      * `volume-attachment`
        * list/create/get/delete
    * Cloud Files (`files`)
      * `account`
        * set/get/update/delete metadata
      * `container`
        * list/create/get/update/delete/empty
        * set/get/update/delete metadata
      * `object`
        * list/upload/get/download/delete
        * set/get/update/delete metadata
        * upload entire directory of files (upload-dir)
    * Cloud Networks (`networks`)
      * `network`
        * list/create/get/update/delete
      * `subnet`
        * list/create/get/update/delete
      * `port`
        * list/create/get/update/delete
      * `security-group`
        * list/create/get/delete
      * `security-group-rule`
        * list/create/get/delete
    * Cloud Block Storage (`block-storage`)
      * `volume`
        * list/create/get/update/delete
      * `snapshot`
        * list/create/get/delete
    * `rack init`: create a `man page` and setup command-completion
    * `rack configure`: interactively create a profile with which to authenticate
    * `rack version`: print out the version and commit hash of the binary
  * Authentication via command-line flags, config file, or environment variables
  * Automatic token caching per service endpoint (with optional `no-cache` flag to disable)
