class NovelmakerObs < Formula
  desc "CLI tool for managing novel projects in Obsidian vaults with OpenAI integration"
  homepage "https://github.com/voilelab/gonovelmaker"
  url "https://github.com/voilelab/gonovelmaker/archive/refs/tags/v0.0.5.tar.gz"
  sha256 "9f58a04686b5733dd1c0934dca36bd69acadffbcbc5867043822fe7805a5a660" # Run: shasum -a 256 <downloaded-file.tar.gz>
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
    system "#{bin}/novelmaker-obs", "--version"
  end
end
