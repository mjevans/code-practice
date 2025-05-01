#!/usr/bin/env python3
# -*- coding: utf_8 -*-
# kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
# https://pycodestyle.pycqa.org/en/latest/intro.html#configuration ignore = E501,W191
# Yes, python authors tend to prefer spaces.  My private codebases prefer tabs for user defined indent.  Tabs should _always_ be preferred, and I'd accept 8, 4 or 2 as reasonable display widths.  Arguably 8 discourages deeply nested code within singule functions.
# Any code going into a shared codebase should (I intend to) follow the lint and coding conventions of that repository.

# 2025 - Michael J Evans
# Code in this file is CC BY-SA 4.0 license https://creativecommons.org/licenses/by-sa/4.0/

# Aside from the tab/space/line length lint defaults and general space sensitivity (and implicit always valid Unicode) of Python: the relative ease of working with JSON and likely XML is far greater than my experience doing the same in Go

import json


def TestJSONs(test_json, key_req=[]):
	p, f, e = 0, 0, 0
	for jsl in test_json:
		try:
			root = json.loads(jsl)
			for k in key_req:
				if k not in root:
					print("Did not find", k, "within", root)
					raise KeyError(k)
					# Python lacks break/continue labels, and also goto to emulate them.
			p += 1
		except KeyError:
			f += 1
		except (BaseException, Exception) as err:
			e += 1
			print(err)
	return p, f, e


def jsonDecTest():
	# test
	key_req = ["id", "name"]
	# pass, pass, errdecode, fail, fail, errdecode?
	# 2 3 2 -- pass pass err p? f? f f err
	test_json = [
		"{\"id\": \"200\", \"name\": \"Test User\"}",
		"\r\n  \t{\"id\": \"200\", \"name\": \"Test User\", \"none\": \"Test User\"}\t\r\n   \n",
		"{id: 200, name: \"Test User\"}",
		"{\"id\": 200, \"name\": \"Test User\"}",
		"{\"id\": \"200\", \"none\": \"Test User\"}",
		"{\"did\": \"200\", \"name\": \"Test User\"}",
		"\"id\": \"200\", \"name\": \"Test User\""
	]

	p, f, e = TestJSONs(test_json, key_req)
	print("\njson decode tests, alternate outcome from encoding/json :\t", (3 == p and 2 == f and 2 == e), "\t results:\t", p, f, e, sep='')
	if 3 == p and 2 == f and 2 == e:
		return
	raise AssertionError("Failed JSON test.")


jsonDecTest()
