package = JSON.parse(File.read(File.join(__dir__, "package.json")))
version = package['version']

Pod::Spec.new do |s|
  s.name             = "TextileIPFS"
  s.version          = version
  s.summary          = package["description"]
  s.requires_arc = true
  s.license      = 'MIT'
  s.homepage     = 'http://www.textile.io'
  s.authors      = { "Aaron Sutula" => "aaron@textile.io" }
  #s.source       = { :git => "https://github.com/textileio/textile-go.git", :tag => 'v#{version}'}
  s.source       = { :git => "https://github.com/textileio/textile-go.git", :branch => 'aaron/rn-npm-module'}
  s.source_files = 'ios/TextileIPFS/*.{h,m}'
  s.preserve_paths  = "**/*.js, **/*.framework"
  s.platform     = :ios, "8.0"
  s.vendored_frameworks = 'Mobile.framework'
end