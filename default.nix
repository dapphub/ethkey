{ stdenv, buildGoPackage, fetchFromGitHub, fetchgit, clang }:

buildGoPackage rec {
  name = "ethkey-${version}";
  version = "0.1";

  goPackagePath = "github.com/dapphub/ethkey";
  hardeningDisable = ["fortify"];
  src = ./.;

  extraSrcs = [
    {
      goPackagePath = "github.com/ethereum/go-ethereum";
      src = fetchFromGitHub {
        owner = "ethereum";
        repo = "go-ethereum";
        rev = "v1.8.1";
        sha256 = "0k7ly9cw68ranksa1fdn7v2lncmlqgabw3qiiyqya2xz3s4aazlf";
      };
    }
    {
      goPackagePath = "gopkg.in/urfave/cli.v1";
      src = fetchFromGitHub {
        owner = "urfave";
        repo = "cli";
        rev = "v1.19.1";
        sha256 = "1ny63c7bfwfrsp7vfkvb4i0xhq4v7yxqnwxa52y4xlfxs4r6v6fg";
      };
    }
  ];

  meta = with stdenv.lib; {
    homepage = https://github.com/dapphub/ethkey;
    description = "Create Ethereum account files";
    license = [licenses.gpl3];
  };
}