class Ignite < Formula
  desc "Build, launch, and maintain any crypto application with Ignite CLI"
  homepage "https://github.com/ignite/cli"
  url "https://github.com/ignite/cli/archive/refs/tags/v28.2.0.tar.gz"
  sha256 "556f953fd7f922354dea64e7b3dade5dd75b3f62ece93167e2ba126cac27602e"
  license "Apache-2.0"

  depends_on "go"
  depends_on "node"

  def install
    system "go", "build", "-mod=readonly", *std_go_args(output: bin/"ignite"), "./ignite/cmd/ignite"
  end

  test do
    ENV["DO_NOT_TRACK"] = "1"
    system bin/"ignite", "s", "chain", "mars"
    assert_predicate testpath/"mars/go.mod", :exist?
  end
end
