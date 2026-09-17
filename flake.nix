{

	inputs = {
		nixpkgs.url = "github:NixOS/nixpkgs/nixos-26.05";
	};

	outputs = { self, nixpkgs, ... }:
		let
			system = "x86_64-linux";
			myOverlays = final: prev: {
				go = prev.go_1_27;
			};
			pkgs = import nixpkgs {
				inherit system;
				overlays = [ myOverlays ];
			};
		in
		{
			devShells.${system}.default = pkgs.mkShell {
				buildInputs = with pkgs; [
					go
				];
			};
		};

}
