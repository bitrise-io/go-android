# Test data files

The files stored here are used in the packages tests:

`empty_file`

This is just an empty file for testing keystore reading with an invalid file.

`pkcs12_type_keystore.jks`

This file is a PKCS12 type keystore and was generated using Android Studio (Build / "Generate Signed Bundle / APK" and going with the create new keystore option).

`jks_type_keystore.keystore`

This file is a JKS type keystore, such a keystore can be generated using the following command:

`keytool -genkey -v -keystore my.keystore -alias my_alias -keyalg RSA -keysize 2048 -validity 1095 -storetype jks -dname "CN=My Common Name,O=My Organisation,C=My Local"`
