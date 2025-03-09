(lang dune 3.18)

(name unit_converter)

(generate_opam_files true)

(source
 (github username/reponame))

(authors "Author Name <author@example.com>")

(maintainers "Maintainer Name <maintainer@example.com>")

(license LICENSE)

(documentation https://url/to/documentation)

(package
 (name unit_converter)
 (synopsis "A short synopsis")
 (description "A longer description")
 (depends core dream dune ocaml lwt_ppx ocaml-lsp-server ocamlformat mlx
  ocamlmerlin-mlx ocamlformat-mlx html_of_jsx)
 (tags
  (topics "to describe" your project)))

(dialect
 (name mlx)
 (implementation
  (extension mlx)
  (merlin_reader mlx)
  (preprocess
   (run mlx-pp %{input-file}))))

; See the complete stanza docs at https://dune.readthedocs.io/en/stable/reference/dune-project/index.html
