# Pomodux Documentation

This directory contains permanent documentation for the Pomodux project, focused on long-term value and development support.

## 📚 Documentation Structure

The repository documentation includes:

1. **[Release Management](release-management.md)** - Release process and approval gates
2. **[Requirements](requirements.md)** - Project requirements and specifications  
3. **[Technical Specifications](technical_specifications.md)** - Technical architecture and design
4. **[Configuration File Specifications](configuration_file_specifications.md)** - Configuration file structure and options
5. **[Development Setup](development-setup.md)** - Development environment and tools
6. **[Go Standards](go-standards.md)** - Go coding standards and conventions
7. **[Logging Standards](logging-standards.md)** - Logging configuration and standards
8. **[Documentation Standards](documentation-standards.md)** - Documentation guidelines and templates
9. **[Plugin Development](plugin-development.md)** - Plugin development guidelines
10. **[CI/CD Pipeline](ci-cd-pipeline.md)** - Comprehensive CI/CD automation and daily operations
11. **[Code Review](code-review.md)** - Code quality assessment and recommendations
12. **[ADR](adr/)** - Architecture Decision Records

## 📁 External Documentation

Planning and temporary documentation is maintained externally to keep the repository focused:

- **Planning & Backlog**: `~/Documents/pomodux/planning/` - Current and future feature planning
- **Implementation Plans**: `~/Documents/pomodux/implementation-plans/` - Detailed feature implementation docs  
- **Historical Records**: `~/Documents/pomodux/releases/` - Historical release documentation
- **Release Retrospectives**: `~/Documents/pomodux/history/` - Release retrospectives and lessons learned

## 🎯 Quick Navigation by Audience

### **For New Contributors**
1. **[Development Setup](development-setup.md)** - Get started with development environment
2. **[Go Standards](go-standards.md)** - Coding standards and practices
3. **[Configuration File Specifications](configuration_file_specifications.md)** - Configuration file structure and options
4. **[Requirements](requirements.md)** - Project requirements and goals

### **For Current Development**
1. **[Release Management](release-management.md)** - Release process and approval gates
2. **[Technical Specifications](technical_specifications.md)** - Technical architecture and design
3. **External Planning**: `~/Documents/pomodux/planning/` - Feature planning and backlog

### **For Historical Reference**
1. **[ADR](adr/)** - Architecture decision records and rationale
2. **External Releases**: `~/Documents/pomodux/releases/` - What was actually delivered in each release
3. **External History**: `~/Documents/pomodux/history/` - Retrospectives and lessons learned

## 📋 Documentation by Purpose

| Document | Purpose | Primary Audience |
|----------|---------|------------------|
| [release-management.md](release-management.md) | Release process and approval gates | Developers, Stakeholders |
| [requirements.md](requirements.md) | Project requirements and specifications | Stakeholders, Developers |
| [technical_specifications.md](technical_specifications.md) | Technical architecture and design | Developers, Architects |
| [configuration_file_specifications.md](configuration_file_specifications.md) | Configuration file structure and options | Users, Developers |
| [development-setup.md](development-setup.md) | Development environment and tools | Developers |
| [go-standards.md](go-standards.md) | Go coding standards and conventions | Developers |
| [logging-standards.md](logging-standards.md) | Logging configuration and standards | Developers, DevOps |
| [documentation-standards.md](documentation-standards.md) | Documentation guidelines and templates | All contributors |
| [plugin-development.md](plugin-development.md) | Plugin development guidelines | Plugin developers |
| [ci-cd-pipeline.md](ci-cd-pipeline.md) | CI/CD automation and daily operations | Developers, DevOps |
| [code-review.md](code-review.md) | Code quality assessment and recommendations | Senior developers, Architects |
| [adr/](adr/) | Architecture Decision Records | Developers, Architects |

## 📊 Current Status

### ✅ Completed Releases
- **Release 0.1.0**: Project Foundation & Core Timer Engine
- **Release 0.2.0**: CLI Interface & Basic Functionality  
- **Release 0.2.1**: Persistent Timer with Keypress Controls
- **Release 0.3.0**: CLI Improvements & Plugin System Foundation

### 📋 Available for Future Releases
- **Terminal User Interface (TUI)** - Interactive terminal interface
- **Plugin System Implementation** - Complete the plugin system foundation
- **Data Export and Import** - Session data management
- **Advanced Statistics** - Enhanced analytics and reporting
- **Logging Enhancements** - Log rotation, analysis tools, enhanced configuration

### 🔗 Key Documentation
- **[Release Management](release-management.md)** - Complete release process and standards
- **External Planning**: `~/Documents/pomodux/planning/` - Feature planning and backlog
- **External Releases**: `~/Documents/pomodux/releases/` - Historical records of completed releases

## 🔄 Documentation Workflow

### **Development Process**
1. **Planning**: Requirements defined in external planning directory
2. **Development**: Work tracked using `.claude/commands/plan-and-implement.md` pattern
3. **Release**: Historical record created in `releases/`
4. **Retrospective**: Lessons learned documented in external history directory

### **Documentation Standards**
- **Permanent Documentation**: Only commit documentation with long-term value
- **Temporary Documentation**: Use external directories for planning and implementation docs
- **Structured Development**: Follow `.claude/commands/plan-and-implement.md` pattern for features
- **Clear Organization**: Separate repository docs from external planning materials

## 🎯 Key Principles

### **Industry Standards Alignment**
- Follows Agile/Scrum principles
- Aligns with DevOps practices
- Supports open source workflows
- Maintains project-specific efficiency

### **User-Centric Organization**
- Easy navigation for different audiences
- Clear purpose for each document type
- Consistent formatting and structure
- Comprehensive coverage of project needs

## 🔗 Related Resources

- **[GitHub Repository](https://github.com/your-org/pomodux)** - Source code and issues
- **[Release Notes](releases/)** - What's new in each release
- **[Contributing Guide](../CONTRIBUTING.md)** - How to contribute
- **[License](../LICENSE)** - Project license

---

**Note**: This streamlined documentation structure follows the principle of keeping permanent documentation in the repository while maintaining planning and temporary documentation externally. This reduces repository complexity while maintaining comprehensive project documentation.

**Last Updated:** 2025-08-09 