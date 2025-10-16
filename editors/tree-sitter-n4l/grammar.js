// Basic Tree-sitter grammar for N4L language (initial draft)
// Focus: sections, contexts (:: markers), aliases (@name), references ($n, $alias), relationships (parenthesized), comments

module.exports = grammar({
  name: "n4l",

  extras: ($) => [/[ \t\f\r\n]/, $.comment],

  rules: {
    source_file: ($) =>
      repeat(
        choice(
          $.section,
          $.context_block,
          $.relation_line,
          $.alias_definition,
          $.reference_line,
          $.todo_block,
          $.statement,
        ),
      ),

    comment: ($) => token(choice(seq("#", /.*/), seq("//", /.*/))),

    section: ($) => seq("-", field("title", /.*/)),

    context_block: ($) =>
      seq(
        field("prefix", optional(choice("+", "-"))),
        "::",
        field("content", /[^:][^\n]*/),
        "::",
      ),

    todo_block: ($) => /[A-Z][A-Z0-9 ]+[A-Z0-9]/,

    alias_definition: ($) => seq("@", field("alias", /[A-Za-z_][A-Za-z0-9_]*/)),

    reference_line: ($) =>
      choice(
        seq("$", field("number_ref", /[0-9]+/)),
        seq("$", field("alias_ref", /[A-Za-z_][A-Za-z0-9_]*/)),
      ),

    relation_line: ($) => seq("(", field("relation", /[^)]+/), ")"),

    statement: ($) => /[^\n]+/,
  },
});
