import argparse, configparser, sys, os

# Written (at home, on my own time and resources) 2025-06-15 mjevans
# version 0.0.0 WIP - general ideas workshopped while brushing up on Python (after years).
# Code in this file LICENSE - CC BY-SA 4.0 license https://creativecommons.org/licenses/by-sa/4.0/
# I also extend a non-revocable, non-exclusive but otherwise unlimited license to my, at the time of writing, contract employer

if sys.version_info >= (3, 5):
    raise AssertionError("Unsupported version of Python, ConfigTool requires 3.5+")

def GetProgName():
    fn = os.path.basename(sys.modules['__main__'].__file__)
    fnext = fn.rfind(".")
    if 0 >= fnext:
        return fn
    return fn[0:fnext]

def DefaultConfigFile(ext = ".conf"):
    return GetProgName() + ext

# XDG path constants

userhome = os.path.expanduser('~')
xdg_data_home = os.path.abspath(os.environ.get('XDG_DATA_HOME') or os.path.join(userhome, '.local', '.share'))
xdg_config_home = os.path.abspath(os.environ.get('XDG_CONFIG_HOME') or os.path.join(userhome, '.config'))
progName = GetProgName()

# print(... , file=sys.stderr) vs os.stderr.write() vs logging
# This is 'early' setup, and a library that interacts with command line arguments which might even specify where config files are.
# Thus 'logging' is premature here, and the file= method probably the most correct / obvious.

class ConfigTool:
    """Take a dict of dicts of configuration defaults to update from ini style config file(s) and command line argparse overrides.  Notably configparser [DEFAULT] behavior is disabled.  Please use [global] in this library for similar use cases."""
    Config = {}
    Types = {}
    ns = None

    def __init__(self, defaults={"global":{}}, types={}):
        self.Config, self.Types = defaults, types
        # self.ReconfigureArgparse() # Delay until just before parsing command line arguments

    def ApplyConfigFiles(searchpaths = [os.path.join(os.sep, "etc", progName), os.path.join(os.sep, "etc", "default", progName), os.path.join(xdg_config_home, progName, DefaultConfigFile()), os.path.join(os.getcwd(), DefaultConfigFile()]):
        minicfg = argparse.ArgumentParser(prog=progName)
        minicfg.add_argument('-v', '--verbose', action='count', default=0)
        minicfg.add_argument('-c', '--config', action='append')
        ns = minicfg.parse_args()
        v = 0 < ns.verbose
        if 'config' in ns:
            searchpaths = ns.config
        for path in searchpaths:
            if os.path.exists(path):
                if v:
                    print("Loading {}".format(v), file=os.stderr)
                self.MergeConfig(v)
            elif v:
                print("Not found, skipping {}".format(v), file=os.stderr)
        self.ReconfigureArgparse()

    def MergeConfig(self, path):
        # configparser looks surprisingly close to the behavior I want, but it doesn't specify exactly how defaults apply in what ordering.
        # more specific settings should __always__ to take precedence over a default, but later configuration files (and command line arguments) should apply first of all unless a more specific setting exists.
        cfg = configparser.ConfigParser(empty_lines_in_values=False, default_section=None)
        try:
            cfg.read(path)
            for sect in cfg.sections():
                # After research online: 3.5+ merged = {**a, **b} while 3.9+ merged = a | b
                self.config[sect] = {**(self.Config.get(sect, {})), **(cfg.[sect])}
        except Exception as err:
            print("FIXME: {}".format(err), file=os.stderr)

    def ReconfigureArgparse(self):
        # This could get really tedious for a large config file...
        # Annoyingly a custom argv parse would be required to support arbitrary
        argp = argparse.ArgumentParser(prog=progName)
        argp.add_argument('-v', '--verbose', action='count', default=0)
        argp.add_argument('-c', '--config', action='append') # repeated here so they get consumed
        argp.add_argument('-s', '--set', action='append', help="--set \"section_key=value\"")
        self.ns = argp.parse_args()
        if 'set' in ns:
            for raw in ns.get("set"): # set is a keyword, but it's also exactly the correct word for the command line
                li = raw.split('=', 1)
                if 2 == len(li):
                    v = li[1].strip()
                    li = li[0].split('_', 1)
                    if 2 == len(li):
                        sect = li[0].strip()
                        k = li[1].strip()
                        if "" != sect && "" != k:
                            if sect in self.Config:
                                self.Config[sect][k] = v
                            else:
                                self.Config[sect] = {k: v}

    def GetConfig():
        return self.Config

    def GetArgNS():
        return self.ns

#

import configparser
cfg = configparser.ConfigParser(empty_lines_in_values=False, default_section=None)
cfg.read_string("[DEFAULT]\ntest=thing\n")
