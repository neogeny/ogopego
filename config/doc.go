// Package config defines configuration types and helpers used across ogopego.
// The key type is Cfg which centralizes compilation and parsing options
// (memoization, whitespace rules, trace settings, etc.). Use DefaultCfg
// to obtain a sensible default and Cfg.Override to merge user-provided
// overrides.
package config
