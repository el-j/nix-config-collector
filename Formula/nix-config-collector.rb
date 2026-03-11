# typed: false
# frozen_string_literal: true

class NixConfigCollector < Formula
  desc "Collect macOS setup and generate nix-darwin + home-manager configuration"
  homepage "https://el-j.github.io/nix-config-collector"
  url "https://github.com/el-j/nix-config-collector/archive/refs/tags/v0.1.0-45.tar.gz"
  sha256 "8846bd3bc7b3701767cba99023af04138b2e6c7a142d5f99511204ad480180cd"
  license "MIT"
  head "https://github.com/el-j/nix-config-collector.git", branch: "main"

  depends_on "go" => :build
  depends_on :macos

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w -X main.version=#{version}"), "./cmd/cli/"
  end

  test do
    assert_match "nix-config-collector version", shell_output("#{bin}/nix-config-collector version")
  end
end
