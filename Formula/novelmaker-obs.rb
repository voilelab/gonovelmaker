class NovelmakerObs < Formula
  desc "CLI tool for managing novel projects in Obsidian vaults with OpenAI integration"
  homepage "https://github.com/voilelab/gonovelmaker"
  url "https://github.com/voilelab/gonovelmaker/archive/refs/tags/v0.0.10.tar.gz"
  sha256 "c88bce97e8b977ac1a1689ca62305cbe35f52695120ac15e1d43d39386891e5b" # Run: shasum -a 256 <downloaded-file.tar.gz>
  license "MIT"
  head "https://github.com/voilelab/gonovelmaker.git", branch: "main"

  depends_on "go" => :build

  def install
    # Build the binary
    cd "cmd/novelmaker-obs" do
      system "go", "build", *std_go_args(ldflags: "-s -w"), "-o", bin/"novelmaker-obs"
    end

    # Install documentation
    doc.install "README.md"
    doc.install "docs" if Dir.exist?("docs")
  end

  test do
    system "#{bin}/novelmaker-obs", "version"
  end
end
