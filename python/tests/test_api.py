"""Integration tests for the ogopego Python API (compile, parse, parse_file).

Requires the bundled Go binary to be present (built via build.py).
"""

import tempfile
from pathlib import Path

import pytest

import ogopego

CALC_GRAMMAR = """\
@@grammar::CALC

start = expression $ ;
expression = | expression '+' term | expression '-' term | term ;
term = | term '*' factor | term '/' factor | factor ;
factor = | '(' expression ')' | number ;
number = /\\d+/ ;
"""


class TestCompile:
    def test_returns_dict(self):
        result = ogopego.compile(CALC_GRAMMAR)
        assert isinstance(result, dict)

    def test_has_name(self):
        result = ogopego.compile(CALC_GRAMMAR)
        assert result["name"] == "CALC"

    def test_has_rules(self):
        result = ogopego.compile(CALC_GRAMMAR)
        assert len(result["rules"]) == 5

    def test_caching_same_object(self):
        a = ogopego.compile(CALC_GRAMMAR)
        b = ogopego.compile(CALC_GRAMMAR)
        assert a is b


class TestParse:
    def test_number(self):
        tree = ogopego.parse(CALC_GRAMMAR, "42")
        assert tree == "42"

    def test_binary_expression(self):
        tree = ogopego.parse(CALC_GRAMMAR, "1+2")
        assert tree == ["1", "+", "2"]

    def test_precedence(self):
        tree = ogopego.parse(CALC_GRAMMAR, "1+2*3")
        assert tree == ["1", "+", ["2", "*", "3"]]

    def test_parentheses(self):
        tree = ogopego.parse(CALC_GRAMMAR, "(1+2)*3")
        assert tree == [["(", ["1", "+", "2"], ")"], "*", "3"]


class TestParseFile:
    def test_matches_parse(self):
        with tempfile.NamedTemporaryFile(mode="w", suffix=".txt", delete=False) as f:
            f.write("1+2")
            path = f.name
        try:
            tree = ogopego.parse_file(CALC_GRAMMAR, path)
            assert tree == ["1", "+", "2"]
        finally:
            Path(path).unlink(missing_ok=True)


class TestOgoError:
    def test_bad_grammar(self):
        with pytest.raises(ogopego.OgoError):
            ogopego.compile("this is not a valid grammar")

    def test_bad_input(self):
        with pytest.raises(ogopego.OgoError):
            ogopego.parse(CALC_GRAMMAR, "this will not parse")
