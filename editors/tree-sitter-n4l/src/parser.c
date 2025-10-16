#include "tree_sitter/parser.h"

#if defined(__GNUC__) || defined(__clang__)
#pragma GCC diagnostic ignored "-Wmissing-field-initializers"
#endif

#define LANGUAGE_VERSION 14
#define STATE_COUNT 22
#define LARGE_STATE_COUNT 4
#define SYMBOL_COUNT 23
#define ALIAS_COUNT 0
#define TOKEN_COUNT 16
#define EXTERNAL_TOKEN_COUNT 0
#define FIELD_COUNT 7
#define MAX_ALIAS_SEQUENCE_LENGTH 4
#define PRODUCTION_ID_COUNT 8

enum ts_symbol_identifiers {
  sym_comment = 1,
  anon_sym_DASH = 2,
  aux_sym_section_token1 = 3,
  anon_sym_PLUS = 4,
  anon_sym_COLON_COLON = 5,
  aux_sym_context_block_token1 = 6,
  sym_todo_block = 7,
  anon_sym_AT = 8,
  aux_sym_alias_definition_token1 = 9,
  anon_sym_DOLLAR = 10,
  aux_sym_reference_line_token1 = 11,
  anon_sym_LPAREN = 12,
  aux_sym_relation_line_token1 = 13,
  anon_sym_RPAREN = 14,
  sym_statement = 15,
  sym_source_file = 16,
  sym_section = 17,
  sym_context_block = 18,
  sym_alias_definition = 19,
  sym_reference_line = 20,
  sym_relation_line = 21,
  aux_sym_source_file_repeat1 = 22,
};

static const char * const ts_symbol_names[] = {
  [ts_builtin_sym_end] = "end",
  [sym_comment] = "comment",
  [anon_sym_DASH] = "-",
  [aux_sym_section_token1] = "section_token1",
  [anon_sym_PLUS] = "+",
  [anon_sym_COLON_COLON] = "::",
  [aux_sym_context_block_token1] = "context_block_token1",
  [sym_todo_block] = "todo_block",
  [anon_sym_AT] = "@",
  [aux_sym_alias_definition_token1] = "alias_definition_token1",
  [anon_sym_DOLLAR] = "$",
  [aux_sym_reference_line_token1] = "reference_line_token1",
  [anon_sym_LPAREN] = "(",
  [aux_sym_relation_line_token1] = "relation_line_token1",
  [anon_sym_RPAREN] = ")",
  [sym_statement] = "statement",
  [sym_source_file] = "source_file",
  [sym_section] = "section",
  [sym_context_block] = "context_block",
  [sym_alias_definition] = "alias_definition",
  [sym_reference_line] = "reference_line",
  [sym_relation_line] = "relation_line",
  [aux_sym_source_file_repeat1] = "source_file_repeat1",
};

static const TSSymbol ts_symbol_map[] = {
  [ts_builtin_sym_end] = ts_builtin_sym_end,
  [sym_comment] = sym_comment,
  [anon_sym_DASH] = anon_sym_DASH,
  [aux_sym_section_token1] = aux_sym_section_token1,
  [anon_sym_PLUS] = anon_sym_PLUS,
  [anon_sym_COLON_COLON] = anon_sym_COLON_COLON,
  [aux_sym_context_block_token1] = aux_sym_context_block_token1,
  [sym_todo_block] = sym_todo_block,
  [anon_sym_AT] = anon_sym_AT,
  [aux_sym_alias_definition_token1] = aux_sym_alias_definition_token1,
  [anon_sym_DOLLAR] = anon_sym_DOLLAR,
  [aux_sym_reference_line_token1] = aux_sym_reference_line_token1,
  [anon_sym_LPAREN] = anon_sym_LPAREN,
  [aux_sym_relation_line_token1] = aux_sym_relation_line_token1,
  [anon_sym_RPAREN] = anon_sym_RPAREN,
  [sym_statement] = sym_statement,
  [sym_source_file] = sym_source_file,
  [sym_section] = sym_section,
  [sym_context_block] = sym_context_block,
  [sym_alias_definition] = sym_alias_definition,
  [sym_reference_line] = sym_reference_line,
  [sym_relation_line] = sym_relation_line,
  [aux_sym_source_file_repeat1] = aux_sym_source_file_repeat1,
};

static const TSSymbolMetadata ts_symbol_metadata[] = {
  [ts_builtin_sym_end] = {
    .visible = false,
    .named = true,
  },
  [sym_comment] = {
    .visible = true,
    .named = true,
  },
  [anon_sym_DASH] = {
    .visible = true,
    .named = false,
  },
  [aux_sym_section_token1] = {
    .visible = false,
    .named = false,
  },
  [anon_sym_PLUS] = {
    .visible = true,
    .named = false,
  },
  [anon_sym_COLON_COLON] = {
    .visible = true,
    .named = false,
  },
  [aux_sym_context_block_token1] = {
    .visible = false,
    .named = false,
  },
  [sym_todo_block] = {
    .visible = true,
    .named = true,
  },
  [anon_sym_AT] = {
    .visible = true,
    .named = false,
  },
  [aux_sym_alias_definition_token1] = {
    .visible = false,
    .named = false,
  },
  [anon_sym_DOLLAR] = {
    .visible = true,
    .named = false,
  },
  [aux_sym_reference_line_token1] = {
    .visible = false,
    .named = false,
  },
  [anon_sym_LPAREN] = {
    .visible = true,
    .named = false,
  },
  [aux_sym_relation_line_token1] = {
    .visible = false,
    .named = false,
  },
  [anon_sym_RPAREN] = {
    .visible = true,
    .named = false,
  },
  [sym_statement] = {
    .visible = true,
    .named = true,
  },
  [sym_source_file] = {
    .visible = true,
    .named = true,
  },
  [sym_section] = {
    .visible = true,
    .named = true,
  },
  [sym_context_block] = {
    .visible = true,
    .named = true,
  },
  [sym_alias_definition] = {
    .visible = true,
    .named = true,
  },
  [sym_reference_line] = {
    .visible = true,
    .named = true,
  },
  [sym_relation_line] = {
    .visible = true,
    .named = true,
  },
  [aux_sym_source_file_repeat1] = {
    .visible = false,
    .named = false,
  },
};

enum ts_field_identifiers {
  field_alias = 1,
  field_alias_ref = 2,
  field_content = 3,
  field_number_ref = 4,
  field_prefix = 5,
  field_relation = 6,
  field_title = 7,
};

static const char * const ts_field_names[] = {
  [0] = NULL,
  [field_alias] = "alias",
  [field_alias_ref] = "alias_ref",
  [field_content] = "content",
  [field_number_ref] = "number_ref",
  [field_prefix] = "prefix",
  [field_relation] = "relation",
  [field_title] = "title",
};

static const TSFieldMapSlice ts_field_map_slices[PRODUCTION_ID_COUNT] = {
  [1] = {.index = 0, .length = 1},
  [2] = {.index = 1, .length = 1},
  [3] = {.index = 2, .length = 1},
  [4] = {.index = 3, .length = 1},
  [5] = {.index = 4, .length = 1},
  [6] = {.index = 5, .length = 1},
  [7] = {.index = 6, .length = 2},
};

static const TSFieldMapEntry ts_field_map_entries[] = {
  [0] =
    {field_title, 1},
  [1] =
    {field_alias, 1},
  [2] =
    {field_alias_ref, 1},
  [3] =
    {field_number_ref, 1},
  [4] =
    {field_content, 1},
  [5] =
    {field_relation, 1},
  [6] =
    {field_content, 2},
    {field_prefix, 0},
};

static const TSSymbol ts_alias_sequences[PRODUCTION_ID_COUNT][MAX_ALIAS_SEQUENCE_LENGTH] = {
  [0] = {0},
};

static const uint16_t ts_non_terminal_alias_map[] = {
  0,
};

static const TSStateId ts_primary_state_ids[STATE_COUNT] = {
  [0] = 0,
  [1] = 1,
  [2] = 2,
  [3] = 3,
  [4] = 4,
  [5] = 5,
  [6] = 6,
  [7] = 7,
  [8] = 8,
  [9] = 9,
  [10] = 10,
  [11] = 11,
  [12] = 12,
  [13] = 13,
  [14] = 14,
  [15] = 15,
  [16] = 16,
  [17] = 17,
  [18] = 18,
  [19] = 19,
  [20] = 20,
  [21] = 21,
};

static bool ts_lex(TSLexer *lexer, TSStateId state) {
  START_LEXER();
  eof = lexer->eof(lexer);
  switch (state) {
    case 0:
      if (eof) ADVANCE(8);
      ADVANCE_MAP(
        '#', 10,
        '$', 33,
        '(', 36,
        ')', 41,
        '+', 17,
        '-', 11,
        '/', 5,
        ':', 6,
        '@', 28,
      );
      if (lookahead == '\t' ||
          lookahead == '\n' ||
          lookahead == '\f' ||
          lookahead == '\r' ||
          lookahead == ' ') SKIP(0);
      if (('0' <= lookahead && lookahead <= '9')) ADVANCE(35);
      if (('A' <= lookahead && lookahead <= 'Z')) ADVANCE(31);
      if (lookahead == '_' ||
          ('a' <= lookahead && lookahead <= 'z')) ADVANCE(32);
      END_STATE();
    case 1:
      if (lookahead == ' ') ADVANCE(1);
      if (('0' <= lookahead && lookahead <= '9') ||
          ('A' <= lookahead && lookahead <= 'Z')) ADVANCE(26);
      END_STATE();
    case 2:
      if (lookahead == '#') ADVANCE(10);
      if (lookahead == '/') ADVANCE(5);
      if (lookahead == '\t' ||
          lookahead == '\n' ||
          lookahead == '\f' ||
          lookahead == '\r' ||
          lookahead == ' ') SKIP(2);
      if (('0' <= lookahead && lookahead <= '9')) ADVANCE(35);
      if (('A' <= lookahead && lookahead <= 'Z') ||
          lookahead == '_' ||
          ('a' <= lookahead && lookahead <= 'z')) ADVANCE(32);
      END_STATE();
    case 3:
      if (lookahead == '#') ADVANCE(10);
      if (lookahead == '/') ADVANCE(23);
      if (lookahead == '\t' ||
          lookahead == '\n' ||
          lookahead == '\f' ||
          lookahead == '\r' ||
          lookahead == ' ') ADVANCE(22);
      if (lookahead != 0 &&
          lookahead != ':') ADVANCE(24);
      END_STATE();
    case 4:
      if (lookahead == '#') ADVANCE(9);
      if (lookahead == '/') ADVANCE(39);
      if (lookahead == '\t' ||
          lookahead == '\n' ||
          lookahead == '\f' ||
          lookahead == '\r' ||
          lookahead == ' ') ADVANCE(38);
      if (lookahead != 0 &&
          lookahead != ')') ADVANCE(40);
      END_STATE();
    case 5:
      if (lookahead == '/') ADVANCE(10);
      END_STATE();
    case 6:
      if (lookahead == ':') ADVANCE(19);
      END_STATE();
    case 7:
      if (eof) ADVANCE(8);
      if (lookahead == '\n') SKIP(7);
      if (lookahead == '#') ADVANCE(10);
      if (lookahead == '$') ADVANCE(34);
      if (lookahead == '(') ADVANCE(37);
      if (lookahead == '+') ADVANCE(18);
      if (lookahead == '-') ADVANCE(12);
      if (lookahead == '/') ADVANCE(44);
      if (lookahead == ':') ADVANCE(45);
      if (lookahead == '@') ADVANCE(29);
      if (lookahead == '\t' ||
          lookahead == '\f' ||
          lookahead == '\r' ||
          lookahead == ' ') ADVANCE(43);
      if (('A' <= lookahead && lookahead <= 'Z')) ADVANCE(46);
      if (lookahead != 0) ADVANCE(47);
      END_STATE();
    case 8:
      ACCEPT_TOKEN(ts_builtin_sym_end);
      END_STATE();
    case 9:
      ACCEPT_TOKEN(sym_comment);
      if (lookahead == '\n') ADVANCE(40);
      if (lookahead == ')') ADVANCE(10);
      if (lookahead != 0) ADVANCE(9);
      END_STATE();
    case 10:
      ACCEPT_TOKEN(sym_comment);
      if (lookahead != 0 &&
          lookahead != '\n') ADVANCE(10);
      END_STATE();
    case 11:
      ACCEPT_TOKEN(anon_sym_DASH);
      END_STATE();
    case 12:
      ACCEPT_TOKEN(anon_sym_DASH);
      if (lookahead != 0 &&
          lookahead != '\n') ADVANCE(47);
      END_STATE();
    case 13:
      ACCEPT_TOKEN(aux_sym_section_token1);
      if (lookahead == '#') ADVANCE(10);
      if (lookahead == '/') ADVANCE(14);
      if (lookahead == ':') ADVANCE(15);
      if (lookahead == '\t' ||
          lookahead == '\f' ||
          lookahead == '\r' ||
          lookahead == ' ') ADVANCE(13);
      if (lookahead != 0 &&
          lookahead != '\t' &&
          lookahead != '\n') ADVANCE(16);
      END_STATE();
    case 14:
      ACCEPT_TOKEN(aux_sym_section_token1);
      if (lookahead == '/') ADVANCE(10);
      if (lookahead != 0 &&
          lookahead != '\n') ADVANCE(16);
      END_STATE();
    case 15:
      ACCEPT_TOKEN(aux_sym_section_token1);
      if (lookahead == ':') ADVANCE(21);
      if (lookahead != 0 &&
          lookahead != '\n') ADVANCE(16);
      END_STATE();
    case 16:
      ACCEPT_TOKEN(aux_sym_section_token1);
      if (lookahead != 0 &&
          lookahead != '\n') ADVANCE(16);
      END_STATE();
    case 17:
      ACCEPT_TOKEN(anon_sym_PLUS);
      END_STATE();
    case 18:
      ACCEPT_TOKEN(anon_sym_PLUS);
      if (lookahead != 0 &&
          lookahead != '\n') ADVANCE(47);
      END_STATE();
    case 19:
      ACCEPT_TOKEN(anon_sym_COLON_COLON);
      END_STATE();
    case 20:
      ACCEPT_TOKEN(anon_sym_COLON_COLON);
      if (lookahead != 0 &&
          lookahead != '\n') ADVANCE(47);
      END_STATE();
    case 21:
      ACCEPT_TOKEN(anon_sym_COLON_COLON);
      if (lookahead != 0 &&
          lookahead != '\n') ADVANCE(16);
      END_STATE();
    case 22:
      ACCEPT_TOKEN(aux_sym_context_block_token1);
      if (lookahead == '\n') ADVANCE(22);
      if (lookahead == '#') ADVANCE(10);
      if (lookahead == '/') ADVANCE(23);
      if (lookahead == ':') ADVANCE(24);
      if (lookahead == '\t' ||
          lookahead == '\f' ||
          lookahead == '\r' ||
          lookahead == ' ') ADVANCE(22);
      if (lookahead != 0) ADVANCE(24);
      END_STATE();
    case 23:
      ACCEPT_TOKEN(aux_sym_context_block_token1);
      if (lookahead == '/') ADVANCE(10);
      if (lookahead != 0 &&
          lookahead != '\n') ADVANCE(24);
      END_STATE();
    case 24:
      ACCEPT_TOKEN(aux_sym_context_block_token1);
      if (lookahead != 0 &&
          lookahead != '\n') ADVANCE(24);
      END_STATE();
    case 25:
      ACCEPT_TOKEN(sym_todo_block);
      if (lookahead == ' ') ADVANCE(1);
      if (lookahead == '_' ||
          ('a' <= lookahead && lookahead <= 'z')) ADVANCE(32);
      if (('0' <= lookahead && lookahead <= '9') ||
          ('A' <= lookahead && lookahead <= 'Z')) ADVANCE(25);
      END_STATE();
    case 26:
      ACCEPT_TOKEN(sym_todo_block);
      if (lookahead == ' ') ADVANCE(1);
      if (('0' <= lookahead && lookahead <= '9') ||
          ('A' <= lookahead && lookahead <= 'Z')) ADVANCE(26);
      END_STATE();
    case 27:
      ACCEPT_TOKEN(sym_todo_block);
      if (lookahead == ' ') ADVANCE(42);
      if (('0' <= lookahead && lookahead <= '9') ||
          ('A' <= lookahead && lookahead <= 'Z')) ADVANCE(27);
      if (lookahead != 0 &&
          lookahead != '\n') ADVANCE(47);
      END_STATE();
    case 28:
      ACCEPT_TOKEN(anon_sym_AT);
      END_STATE();
    case 29:
      ACCEPT_TOKEN(anon_sym_AT);
      if (lookahead != 0 &&
          lookahead != '\n') ADVANCE(47);
      END_STATE();
    case 30:
      ACCEPT_TOKEN(aux_sym_alias_definition_token1);
      if (lookahead == ' ') ADVANCE(1);
      if (lookahead == '_' ||
          ('a' <= lookahead && lookahead <= 'z')) ADVANCE(32);
      if (('0' <= lookahead && lookahead <= '9') ||
          ('A' <= lookahead && lookahead <= 'Z')) ADVANCE(25);
      END_STATE();
    case 31:
      ACCEPT_TOKEN(aux_sym_alias_definition_token1);
      if (lookahead == ' ') ADVANCE(1);
      if (lookahead == '_' ||
          ('a' <= lookahead && lookahead <= 'z')) ADVANCE(32);
      if (('0' <= lookahead && lookahead <= '9') ||
          ('A' <= lookahead && lookahead <= 'Z')) ADVANCE(30);
      END_STATE();
    case 32:
      ACCEPT_TOKEN(aux_sym_alias_definition_token1);
      if (('0' <= lookahead && lookahead <= '9') ||
          ('A' <= lookahead && lookahead <= 'Z') ||
          lookahead == '_' ||
          ('a' <= lookahead && lookahead <= 'z')) ADVANCE(32);
      END_STATE();
    case 33:
      ACCEPT_TOKEN(anon_sym_DOLLAR);
      END_STATE();
    case 34:
      ACCEPT_TOKEN(anon_sym_DOLLAR);
      if (lookahead != 0 &&
          lookahead != '\n') ADVANCE(47);
      END_STATE();
    case 35:
      ACCEPT_TOKEN(aux_sym_reference_line_token1);
      if (('0' <= lookahead && lookahead <= '9')) ADVANCE(35);
      END_STATE();
    case 36:
      ACCEPT_TOKEN(anon_sym_LPAREN);
      END_STATE();
    case 37:
      ACCEPT_TOKEN(anon_sym_LPAREN);
      if (lookahead != 0 &&
          lookahead != '\n') ADVANCE(47);
      END_STATE();
    case 38:
      ACCEPT_TOKEN(aux_sym_relation_line_token1);
      if (lookahead == '#') ADVANCE(9);
      if (lookahead == '/') ADVANCE(39);
      if (lookahead == '\t' ||
          lookahead == '\n' ||
          lookahead == '\f' ||
          lookahead == '\r' ||
          lookahead == ' ') ADVANCE(38);
      if (lookahead != 0 &&
          lookahead != ')') ADVANCE(40);
      END_STATE();
    case 39:
      ACCEPT_TOKEN(aux_sym_relation_line_token1);
      if (lookahead == '/') ADVANCE(9);
      if (lookahead != 0 &&
          lookahead != ')') ADVANCE(40);
      END_STATE();
    case 40:
      ACCEPT_TOKEN(aux_sym_relation_line_token1);
      if (lookahead != 0 &&
          lookahead != ')') ADVANCE(40);
      END_STATE();
    case 41:
      ACCEPT_TOKEN(anon_sym_RPAREN);
      END_STATE();
    case 42:
      ACCEPT_TOKEN(sym_statement);
      if (lookahead == ' ') ADVANCE(42);
      if (('0' <= lookahead && lookahead <= '9') ||
          ('A' <= lookahead && lookahead <= 'Z')) ADVANCE(27);
      if (lookahead != 0 &&
          lookahead != '\n') ADVANCE(47);
      END_STATE();
    case 43:
      ACCEPT_TOKEN(sym_statement);
      ADVANCE_MAP(
        '#', 10,
        '$', 34,
        '(', 37,
        '+', 18,
        '-', 12,
        '/', 44,
        ':', 45,
        '@', 29,
        '\t', 43,
        '\f', 43,
        '\r', 43,
        ' ', 43,
      );
      if (('A' <= lookahead && lookahead <= 'Z')) ADVANCE(46);
      if (lookahead != 0 &&
          lookahead != '\t' &&
          lookahead != '\n') ADVANCE(47);
      END_STATE();
    case 44:
      ACCEPT_TOKEN(sym_statement);
      if (lookahead == '/') ADVANCE(10);
      if (lookahead != 0 &&
          lookahead != '\n') ADVANCE(47);
      END_STATE();
    case 45:
      ACCEPT_TOKEN(sym_statement);
      if (lookahead == ':') ADVANCE(20);
      if (lookahead != 0 &&
          lookahead != '\n') ADVANCE(47);
      END_STATE();
    case 46:
      ACCEPT_TOKEN(sym_statement);
      if (lookahead == ' ' ||
          ('0' <= lookahead && lookahead <= '9') ||
          ('A' <= lookahead && lookahead <= 'Z')) ADVANCE(42);
      if (lookahead != 0 &&
          lookahead != '\n') ADVANCE(47);
      END_STATE();
    case 47:
      ACCEPT_TOKEN(sym_statement);
      if (lookahead != 0 &&
          lookahead != '\n') ADVANCE(47);
      END_STATE();
    default:
      return false;
  }
}

static const TSLexMode ts_lex_modes[STATE_COUNT] = {
  [0] = {.lex_state = 0},
  [1] = {.lex_state = 7},
  [2] = {.lex_state = 7},
  [3] = {.lex_state = 7},
  [4] = {.lex_state = 7},
  [5] = {.lex_state = 7},
  [6] = {.lex_state = 7},
  [7] = {.lex_state = 7},
  [8] = {.lex_state = 7},
  [9] = {.lex_state = 7},
  [10] = {.lex_state = 7},
  [11] = {.lex_state = 2},
  [12] = {.lex_state = 13},
  [13] = {.lex_state = 3},
  [14] = {.lex_state = 2},
  [15] = {.lex_state = 4},
  [16] = {.lex_state = 0},
  [17] = {.lex_state = 3},
  [18] = {.lex_state = 0},
  [19] = {.lex_state = 0},
  [20] = {.lex_state = 0},
  [21] = {.lex_state = 0},
};

static const uint16_t ts_parse_table[LARGE_STATE_COUNT][SYMBOL_COUNT] = {
  [0] = {
    [ts_builtin_sym_end] = ACTIONS(1),
    [sym_comment] = ACTIONS(3),
    [anon_sym_DASH] = ACTIONS(1),
    [anon_sym_PLUS] = ACTIONS(1),
    [anon_sym_COLON_COLON] = ACTIONS(1),
    [sym_todo_block] = ACTIONS(1),
    [anon_sym_AT] = ACTIONS(1),
    [aux_sym_alias_definition_token1] = ACTIONS(1),
    [anon_sym_DOLLAR] = ACTIONS(1),
    [aux_sym_reference_line_token1] = ACTIONS(1),
    [anon_sym_LPAREN] = ACTIONS(1),
    [anon_sym_RPAREN] = ACTIONS(1),
  },
  [1] = {
    [sym_source_file] = STATE(16),
    [sym_section] = STATE(2),
    [sym_context_block] = STATE(2),
    [sym_alias_definition] = STATE(2),
    [sym_reference_line] = STATE(2),
    [sym_relation_line] = STATE(2),
    [aux_sym_source_file_repeat1] = STATE(2),
    [ts_builtin_sym_end] = ACTIONS(5),
    [sym_comment] = ACTIONS(7),
    [anon_sym_DASH] = ACTIONS(9),
    [anon_sym_PLUS] = ACTIONS(11),
    [anon_sym_COLON_COLON] = ACTIONS(13),
    [sym_todo_block] = ACTIONS(15),
    [anon_sym_AT] = ACTIONS(17),
    [anon_sym_DOLLAR] = ACTIONS(19),
    [anon_sym_LPAREN] = ACTIONS(21),
    [sym_statement] = ACTIONS(15),
  },
  [2] = {
    [sym_section] = STATE(3),
    [sym_context_block] = STATE(3),
    [sym_alias_definition] = STATE(3),
    [sym_reference_line] = STATE(3),
    [sym_relation_line] = STATE(3),
    [aux_sym_source_file_repeat1] = STATE(3),
    [ts_builtin_sym_end] = ACTIONS(23),
    [sym_comment] = ACTIONS(7),
    [anon_sym_DASH] = ACTIONS(9),
    [anon_sym_PLUS] = ACTIONS(11),
    [anon_sym_COLON_COLON] = ACTIONS(13),
    [sym_todo_block] = ACTIONS(25),
    [anon_sym_AT] = ACTIONS(17),
    [anon_sym_DOLLAR] = ACTIONS(19),
    [anon_sym_LPAREN] = ACTIONS(21),
    [sym_statement] = ACTIONS(25),
  },
  [3] = {
    [sym_section] = STATE(3),
    [sym_context_block] = STATE(3),
    [sym_alias_definition] = STATE(3),
    [sym_reference_line] = STATE(3),
    [sym_relation_line] = STATE(3),
    [aux_sym_source_file_repeat1] = STATE(3),
    [ts_builtin_sym_end] = ACTIONS(27),
    [sym_comment] = ACTIONS(7),
    [anon_sym_DASH] = ACTIONS(29),
    [anon_sym_PLUS] = ACTIONS(32),
    [anon_sym_COLON_COLON] = ACTIONS(35),
    [sym_todo_block] = ACTIONS(38),
    [anon_sym_AT] = ACTIONS(41),
    [anon_sym_DOLLAR] = ACTIONS(44),
    [anon_sym_LPAREN] = ACTIONS(47),
    [sym_statement] = ACTIONS(38),
  },
};

static const uint16_t ts_small_parse_table[] = {
  [0] = 3,
    ACTIONS(7), 1,
      sym_comment,
    ACTIONS(50), 1,
      ts_builtin_sym_end,
    ACTIONS(52), 8,
      anon_sym_DASH,
      anon_sym_PLUS,
      anon_sym_COLON_COLON,
      sym_todo_block,
      anon_sym_AT,
      anon_sym_DOLLAR,
      anon_sym_LPAREN,
      sym_statement,
  [17] = 3,
    ACTIONS(7), 1,
      sym_comment,
    ACTIONS(54), 1,
      ts_builtin_sym_end,
    ACTIONS(56), 8,
      anon_sym_DASH,
      anon_sym_PLUS,
      anon_sym_COLON_COLON,
      sym_todo_block,
      anon_sym_AT,
      anon_sym_DOLLAR,
      anon_sym_LPAREN,
      sym_statement,
  [34] = 3,
    ACTIONS(7), 1,
      sym_comment,
    ACTIONS(58), 1,
      ts_builtin_sym_end,
    ACTIONS(60), 8,
      anon_sym_DASH,
      anon_sym_PLUS,
      anon_sym_COLON_COLON,
      sym_todo_block,
      anon_sym_AT,
      anon_sym_DOLLAR,
      anon_sym_LPAREN,
      sym_statement,
  [51] = 3,
    ACTIONS(7), 1,
      sym_comment,
    ACTIONS(62), 1,
      ts_builtin_sym_end,
    ACTIONS(64), 8,
      anon_sym_DASH,
      anon_sym_PLUS,
      anon_sym_COLON_COLON,
      sym_todo_block,
      anon_sym_AT,
      anon_sym_DOLLAR,
      anon_sym_LPAREN,
      sym_statement,
  [68] = 3,
    ACTIONS(7), 1,
      sym_comment,
    ACTIONS(66), 1,
      ts_builtin_sym_end,
    ACTIONS(68), 8,
      anon_sym_DASH,
      anon_sym_PLUS,
      anon_sym_COLON_COLON,
      sym_todo_block,
      anon_sym_AT,
      anon_sym_DOLLAR,
      anon_sym_LPAREN,
      sym_statement,
  [85] = 3,
    ACTIONS(7), 1,
      sym_comment,
    ACTIONS(70), 1,
      ts_builtin_sym_end,
    ACTIONS(72), 8,
      anon_sym_DASH,
      anon_sym_PLUS,
      anon_sym_COLON_COLON,
      sym_todo_block,
      anon_sym_AT,
      anon_sym_DOLLAR,
      anon_sym_LPAREN,
      sym_statement,
  [102] = 3,
    ACTIONS(7), 1,
      sym_comment,
    ACTIONS(74), 1,
      ts_builtin_sym_end,
    ACTIONS(76), 8,
      anon_sym_DASH,
      anon_sym_PLUS,
      anon_sym_COLON_COLON,
      sym_todo_block,
      anon_sym_AT,
      anon_sym_DOLLAR,
      anon_sym_LPAREN,
      sym_statement,
  [119] = 3,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(78), 1,
      aux_sym_alias_definition_token1,
    ACTIONS(80), 1,
      aux_sym_reference_line_token1,
  [129] = 3,
    ACTIONS(7), 1,
      sym_comment,
    ACTIONS(82), 1,
      aux_sym_section_token1,
    ACTIONS(84), 1,
      anon_sym_COLON_COLON,
  [139] = 2,
    ACTIONS(7), 1,
      sym_comment,
    ACTIONS(86), 1,
      aux_sym_context_block_token1,
  [146] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(88), 1,
      aux_sym_alias_definition_token1,
  [153] = 2,
    ACTIONS(7), 1,
      sym_comment,
    ACTIONS(90), 1,
      aux_sym_relation_line_token1,
  [160] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(92), 1,
      ts_builtin_sym_end,
  [167] = 2,
    ACTIONS(7), 1,
      sym_comment,
    ACTIONS(94), 1,
      aux_sym_context_block_token1,
  [174] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(96), 1,
      anon_sym_COLON_COLON,
  [181] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(98), 1,
      anon_sym_RPAREN,
  [188] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(100), 1,
      anon_sym_COLON_COLON,
  [195] = 2,
    ACTIONS(3), 1,
      sym_comment,
    ACTIONS(102), 1,
      anon_sym_COLON_COLON,
};

static const uint32_t ts_small_parse_table_map[] = {
  [SMALL_STATE(4)] = 0,
  [SMALL_STATE(5)] = 17,
  [SMALL_STATE(6)] = 34,
  [SMALL_STATE(7)] = 51,
  [SMALL_STATE(8)] = 68,
  [SMALL_STATE(9)] = 85,
  [SMALL_STATE(10)] = 102,
  [SMALL_STATE(11)] = 119,
  [SMALL_STATE(12)] = 129,
  [SMALL_STATE(13)] = 139,
  [SMALL_STATE(14)] = 146,
  [SMALL_STATE(15)] = 153,
  [SMALL_STATE(16)] = 160,
  [SMALL_STATE(17)] = 167,
  [SMALL_STATE(18)] = 174,
  [SMALL_STATE(19)] = 181,
  [SMALL_STATE(20)] = 188,
  [SMALL_STATE(21)] = 195,
};

static const TSParseActionEntry ts_parse_actions[] = {
  [0] = {.entry = {.count = 0, .reusable = false}},
  [1] = {.entry = {.count = 1, .reusable = false}}, RECOVER(),
  [3] = {.entry = {.count = 1, .reusable = true}}, SHIFT_EXTRA(),
  [5] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_source_file, 0, 0, 0),
  [7] = {.entry = {.count = 1, .reusable = false}}, SHIFT_EXTRA(),
  [9] = {.entry = {.count = 1, .reusable = false}}, SHIFT(12),
  [11] = {.entry = {.count = 1, .reusable = false}}, SHIFT(20),
  [13] = {.entry = {.count = 1, .reusable = false}}, SHIFT(13),
  [15] = {.entry = {.count = 1, .reusable = false}}, SHIFT(2),
  [17] = {.entry = {.count = 1, .reusable = false}}, SHIFT(14),
  [19] = {.entry = {.count = 1, .reusable = false}}, SHIFT(11),
  [21] = {.entry = {.count = 1, .reusable = false}}, SHIFT(15),
  [23] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_source_file, 1, 0, 0),
  [25] = {.entry = {.count = 1, .reusable = false}}, SHIFT(3),
  [27] = {.entry = {.count = 1, .reusable = true}}, REDUCE(aux_sym_source_file_repeat1, 2, 0, 0),
  [29] = {.entry = {.count = 2, .reusable = false}}, REDUCE(aux_sym_source_file_repeat1, 2, 0, 0), SHIFT_REPEAT(12),
  [32] = {.entry = {.count = 2, .reusable = false}}, REDUCE(aux_sym_source_file_repeat1, 2, 0, 0), SHIFT_REPEAT(20),
  [35] = {.entry = {.count = 2, .reusable = false}}, REDUCE(aux_sym_source_file_repeat1, 2, 0, 0), SHIFT_REPEAT(13),
  [38] = {.entry = {.count = 2, .reusable = false}}, REDUCE(aux_sym_source_file_repeat1, 2, 0, 0), SHIFT_REPEAT(3),
  [41] = {.entry = {.count = 2, .reusable = false}}, REDUCE(aux_sym_source_file_repeat1, 2, 0, 0), SHIFT_REPEAT(14),
  [44] = {.entry = {.count = 2, .reusable = false}}, REDUCE(aux_sym_source_file_repeat1, 2, 0, 0), SHIFT_REPEAT(11),
  [47] = {.entry = {.count = 2, .reusable = false}}, REDUCE(aux_sym_source_file_repeat1, 2, 0, 0), SHIFT_REPEAT(15),
  [50] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_reference_line, 2, 0, 4),
  [52] = {.entry = {.count = 1, .reusable = false}}, REDUCE(sym_reference_line, 2, 0, 4),
  [54] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_section, 2, 0, 1),
  [56] = {.entry = {.count = 1, .reusable = false}}, REDUCE(sym_section, 2, 0, 1),
  [58] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_alias_definition, 2, 0, 2),
  [60] = {.entry = {.count = 1, .reusable = false}}, REDUCE(sym_alias_definition, 2, 0, 2),
  [62] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_reference_line, 2, 0, 3),
  [64] = {.entry = {.count = 1, .reusable = false}}, REDUCE(sym_reference_line, 2, 0, 3),
  [66] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_context_block, 3, 0, 5),
  [68] = {.entry = {.count = 1, .reusable = false}}, REDUCE(sym_context_block, 3, 0, 5),
  [70] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_relation_line, 3, 0, 6),
  [72] = {.entry = {.count = 1, .reusable = false}}, REDUCE(sym_relation_line, 3, 0, 6),
  [74] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_context_block, 4, 0, 7),
  [76] = {.entry = {.count = 1, .reusable = false}}, REDUCE(sym_context_block, 4, 0, 7),
  [78] = {.entry = {.count = 1, .reusable = true}}, SHIFT(7),
  [80] = {.entry = {.count = 1, .reusable = true}}, SHIFT(4),
  [82] = {.entry = {.count = 1, .reusable = false}}, SHIFT(5),
  [84] = {.entry = {.count = 1, .reusable = false}}, SHIFT(17),
  [86] = {.entry = {.count = 1, .reusable = false}}, SHIFT(18),
  [88] = {.entry = {.count = 1, .reusable = true}}, SHIFT(6),
  [90] = {.entry = {.count = 1, .reusable = false}}, SHIFT(19),
  [92] = {.entry = {.count = 1, .reusable = true}},  ACCEPT_INPUT(),
  [94] = {.entry = {.count = 1, .reusable = false}}, SHIFT(21),
  [96] = {.entry = {.count = 1, .reusable = true}}, SHIFT(8),
  [98] = {.entry = {.count = 1, .reusable = true}}, SHIFT(9),
  [100] = {.entry = {.count = 1, .reusable = true}}, SHIFT(17),
  [102] = {.entry = {.count = 1, .reusable = true}}, SHIFT(10),
};

#ifdef __cplusplus
extern "C" {
#endif
#ifdef TREE_SITTER_HIDE_SYMBOLS
#define TS_PUBLIC
#elif defined(_WIN32)
#define TS_PUBLIC __declspec(dllexport)
#else
#define TS_PUBLIC __attribute__((visibility("default")))
#endif

TS_PUBLIC const TSLanguage *tree_sitter_n4l(void) {
  static const TSLanguage language = {
    .version = LANGUAGE_VERSION,
    .symbol_count = SYMBOL_COUNT,
    .alias_count = ALIAS_COUNT,
    .token_count = TOKEN_COUNT,
    .external_token_count = EXTERNAL_TOKEN_COUNT,
    .state_count = STATE_COUNT,
    .large_state_count = LARGE_STATE_COUNT,
    .production_id_count = PRODUCTION_ID_COUNT,
    .field_count = FIELD_COUNT,
    .max_alias_sequence_length = MAX_ALIAS_SEQUENCE_LENGTH,
    .parse_table = &ts_parse_table[0][0],
    .small_parse_table = ts_small_parse_table,
    .small_parse_table_map = ts_small_parse_table_map,
    .parse_actions = ts_parse_actions,
    .symbol_names = ts_symbol_names,
    .field_names = ts_field_names,
    .field_map_slices = ts_field_map_slices,
    .field_map_entries = ts_field_map_entries,
    .symbol_metadata = ts_symbol_metadata,
    .public_symbol_map = ts_symbol_map,
    .alias_map = ts_non_terminal_alias_map,
    .alias_sequences = &ts_alias_sequences[0][0],
    .lex_modes = ts_lex_modes,
    .lex_fn = ts_lex,
    .primary_state_ids = ts_primary_state_ids,
  };
  return &language;
}
#ifdef __cplusplus
}
#endif
