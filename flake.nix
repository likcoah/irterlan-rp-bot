{

	inputs = {
		nixpkgs.url = "github:NixOS/nixpkgs/nixos-26.05";
	};

	outputs = { self, nixpkgs, ... }:
		let
			system = "x86_64-linux";
			pkgs = nixpkgs.legacyPackages.${system};
		in
		{
			devShells.${system}.default = pkgs.mkShell {
				buildInputs = with pkgs; [
					go
				];
			};
		};

}
