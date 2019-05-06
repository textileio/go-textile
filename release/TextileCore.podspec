Pod::Spec.new do |spec|
  spec.name         = "TextileCore"
  spec.version      = "<version>"
  spec.summary      = "Encrypted, recoverable, schema-based, cross-application data storage built on IPFS and LibP2P"
  spec.description  = <<-DESC
                      Objective C framework and Protobuf files generated from go-textile. You should
                      not usually use this pod directly, but instead use the Textile pod.
                    DESC
  spec.homepage     = "https://github.com/textileio/go-textile"
  spec.license      = "MIT"
  spec.author       = { "textile.io" => "contact@textile.io" }
  spec.platform     = :ios, "7.0"
  spec.source       = spec.source = { :http => 'https://github.com/textileio/go-textile/releases/download/v<version>/go-textile_v<version>_ios-framework.tar.gz' }
  spec.source_files = "protos"
  spec.vendored_frameworks = 'Mobile.framework'
  spec.requires_arc = false
  spec.dependency "Protobuf", "~> 3.7"
  spec.pod_target_xcconfig = { 'GCC_PREPROCESSOR_DEFINITIONS' => '$(inherited) GPB_USE_PROTOBUF_FRAMEWORK_IMPORTS=1', 'OTHER_LDFLAGS[arch=i386]' => '-Wl,-read_only_relocs,suppress' }
end
