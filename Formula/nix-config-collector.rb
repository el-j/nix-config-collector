# typed: false
# frozen_string_literal: true

class NixConfigCollector < Formula
  desc "Collect macOS setup and generate nix-darwin + home-manager configuration"
  homepage "https://el-j.github.io/nix-config-collector"
  url "https://github.com/el-j/nix-config-collector/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "PLACEHOLDER_UPDATE_SHA256_ON_RELEASE" # TODO: replace with `sha256sum v0.1.0.tar.gz` value at release time
  license "MIT"
  head "https://github.com/el-j/nix-config-collector.git", branch: "main"

  depends_on "go" => :build
  depends_on :macos

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "./cmd/cli/"
  end

  test do
    assert_match "nix-config-collector version", shell_output("#{bin}/nix-config-collector version")
  end
end
