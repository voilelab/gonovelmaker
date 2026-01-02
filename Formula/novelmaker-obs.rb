class NovelmakerObs < Formula
  desc "CLI tool for managing novel projects in Obsidian vaults with OpenAI integration"
  homepage "https://github.com/voilelab/gonovelmaker"
  url "https://github.com/voilelab/gonovelmaker/archive/refs/tags/v0.0.6.tar.gz"
  sha256 "5820006244a90097f8a090eedf976dfa1f09b8eaeaeb88c5b62630a4aa3f79a1" # Run: shasum -a 256 <downloaded-file.tar.gz>
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
