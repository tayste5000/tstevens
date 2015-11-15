from Bio.SeqUtils.ProtParam import ProteinAnalysis
from Bio.SeqUtils import ProtParamData
import sys
import json

inp = json.loads(sys.argv[1])

seq = inp["Sequence"]

X = ProteinAnalysis(seq)

data = dict()

if "MW" in inp["Options"]:
	data["MW"] = X.molecular_weight()

if "EC280" in inp["Options"]:
	aa_count = X.count_amino_acids()
	if "hasDisulfide" in inp["Options"]:
		data["EC280"] = 1490 * aa_count["Y"] + 5500 * aa_count["W"] + 62.5 * aa_count["C"]
	else:
		data["EC280"] = 1490 * aa_count["Y"] + 5500 * aa_count["W"]

if "PI" in inp["Options"]:
	data["PI"] = X.isoelectric_point()

print json.dumps(data)