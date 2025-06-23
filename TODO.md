# Security TODO

## Path Traversal Protection

- [ ] Add path validation function to prevent `../` traversal attacks
- [ ] Resolve all paths with `filepath.Clean()` and `filepath.Abs()`
- [ ] Check that resolved paths stay within working directory boundary
- [ ] Reject absolute paths and paths containing `..` segments
- [ ] Apply validation to all file tools: `ReadFile`, `ListFiles`, `EditFile`

## Additional Security Considerations

- [ ] Consider implementing file size limits for reads/writes
- [ ] Add file type restrictions (e.g., no binary files, executables)
- [ ] Implement rate limiting for file operations
- [ ] Add logging for all file system access attempts
- [ ] Consider running in containerized environment for additional isolation