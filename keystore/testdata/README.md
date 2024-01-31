# Test data files

The files stored here are used in the packages tests:

`empty_file`

This is just an empty file for testing keystore reading with an invalid file.

`<keystore_name>.pkcs12`

These files are PKCS12 type keystores and were generated using Android Studio (Build / "Generate Signed Bundle / APK" and going with the create new keystore option).

`<keystore_name>.jks`

These files are JKS type keystores, such a keystore can be generated using keytool:

`keytool -genkey -v -keystore my.keystore -alias my_alias -keyalg RSA -keysize 2048 -validity 1095 -storetype jks -dname "CN=My Common Name,O=My Organisation,C=My Local"`
