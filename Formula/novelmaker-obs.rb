class NovelmakerObs < Formula
  desc "CLI tool for managing novel projects in Obsidian vaults with OpenAI integration"
  homepage "https://github.com/voilelab/gonovelmaker"
  url "https://github.com/voilelab/gonovelmaker/archive/refs/tags/v0.0.7.tar.gz"
  sha256 "61083cbf31fd3435f4ee62b94704e60a22094b1b430231cfb8a4de4d183e66a3" # Run: shasum -a 256 <downloaded-file.tar.gz>
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
