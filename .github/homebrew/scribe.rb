class Scribe < Formula
  desc "A lightweight static site generator built in pure Go with zero external dependencies"
  homepage "https://github.com/dikaio/scribe"
  url "https://github.com/dikaio/scribe/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "TO_BE_REPLACED_AFTER_FIRST_RELEASE"
  license "MIT"
  head "https://github.com/dikaio/scribe.git", branch: "main"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "./cmd/scribe"
  end

  test do
    mkdir_p "mysite"
    cd "mysite" do
      system bin/"scribe", "new", "site", "."
      assert_predicate testpath/"mysite/config.jsonc", :exist?
    end
  end
end