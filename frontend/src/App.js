import React, { useState, useCallback } from 'react';
import { useDropzone } from 'react-dropzone';
import { Upload, Download, Settings, Code, FileText, Zap, Github, Terminal } from 'lucide-react';

function App() {
  const [csvData, setCsvData] = useState('');
  const [jsonOutput, setJsonOutput] = useState('');
  const [loading, setLoading] = useState(false);
  const [options, setOptions] = useState({
    delimiter: ',',
    has_header: true,
    output_format: 'array',
    pretty_print: true,
    infer_types: true
  });
  const [showOptions, setShowOptions] = useState(false);

  const onDrop = useCallback((acceptedFiles) => {
    const file = acceptedFiles[0];
    if (file) {
      const reader = new FileReader();
      reader.onload = (e) => {
        setCsvData(e.target.result);
      };
      reader.readAsText(file);
    }
  }, []);

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop,
    accept: {
      'text/csv': ['.csv'],
      'text/plain': ['.txt']
    },
    multiple: false
  });

  const convertCSV = async () => {
    if (!csvData.trim()) return;
    
    setLoading(true);
    try {
      const response = await fetch('/convert', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          csv_data: csvData,
          options: options
        })
      });

      const result = await response.json();
      if (result.success) {
        setJsonOutput(JSON.stringify(result.data, null, options.pretty_print ? 2 : 0));
      } else {
        setJsonOutput(`Error: ${result.error}`);
      }
    } catch (error) {
      setJsonOutput(`Error: ${error.message}`);
    }
    setLoading(false);
  };

  const downloadJSON = () => {
    if (!jsonOutput) return;
    
    const blob = new Blob([jsonOutput], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'converted.json';
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  };

  const sampleCSV = `name,age,city,active
John Doe,30,New York,true
Jane Smith,25,Los Angeles,false
Bob Johnson,35,Chicago,true`;

  const loadSample = () => {
    setCsvData(sampleCSV);
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-dark-bg via-dark-bg to-dark-card">
      {/* Header */}
      <header className="border-b border-dark-border bg-dark-card/50 backdrop-blur-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-3">
              <div className="w-10 h-10 bg-gradient-to-r from-cyber-blue to-cyber-purple rounded-lg flex items-center justify-center animate-pulse-slow">
                <Zap className="w-6 h-6 text-white" />
              </div>
              <div>
                <h1 className="text-2xl font-bold bg-gradient-to-r from-cyber-blue to-cyber-purple bg-clip-text text-transparent">
                  CSV2JSON
                </h1>
                <p className="text-sm text-gray-400">Modern converter with CLI, API & Web UI</p>
              </div>
            </div>
            <div className="flex items-center space-x-4">
              <a href="https://github.com" className="text-gray-400 hover:text-cyber-blue transition-colors">
                <Github className="w-5 h-5" />
              </a>
              <div className="flex items-center space-x-2 text-sm text-gray-400">
                <Terminal className="w-4 h-4" />
                <span>CLI Available</span>
              </div>
            </div>
          </div>
        </div>
      </header>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Hero Section */}
        <div className="text-center mb-12">
          <h2 className="text-4xl font-bold mb-4 bg-gradient-to-r from-cyber-blue via-cyber-purple to-cyber-pink bg-clip-text text-transparent animate-float">
            Transform CSV to JSON
          </h2>
          <p className="text-xl text-gray-300 mb-8">
            Lightning-fast conversion with advanced options and beautiful output
          </p>
          
          {/* Quick Stats */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
            <div className="cyber-border p-6 text-center">
              <FileText className="w-8 h-8 text-cyber-blue mx-auto mb-2" />
              <div className="text-2xl font-bold text-cyber-blue">CSV</div>
              <div className="text-sm text-gray-400">Input Format</div>
            </div>
            <div className="cyber-border p-6 text-center animate-glow">
              <Zap className="w-8 h-8 text-cyber-purple mx-auto mb-2" />
              <div className="text-2xl font-bold text-cyber-purple">Fast</div>
              <div className="text-sm text-gray-400">Conversion</div>
            </div>
            <div className="cyber-border p-6 text-center">
              <Code className="w-8 h-8 text-cyber-pink mx-auto mb-2" />
              <div className="text-2xl font-bold text-cyber-pink">JSON</div>
              <div className="text-sm text-gray-400">Output Format</div>
            </div>
          </div>
        </div>

        {/* Main Interface */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          {/* Input Section */}
          <div className="space-y-6">
            <div className="flex items-center justify-between">
              <h3 className="text-xl font-semibold text-cyber-blue glow-text">CSV Input</h3>
              <div className="flex space-x-2">
                <button
                  onClick={loadSample}
                  className="px-3 py-1 text-sm bg-dark-card border border-cyber-blue/30 rounded hover:border-cyber-blue transition-colors"
                >
                  Load Sample
                </button>
                <button
                  onClick={() => setShowOptions(!showOptions)}
                  className="p-2 bg-dark-card border border-cyber-blue/30 rounded hover:border-cyber-blue transition-colors"
                >
                  <Settings className="w-4 h-4" />
                </button>
              </div>
            </div>

            {/* File Drop Zone */}
            <div
              {...getRootProps()}
              className={`cyber-border p-8 text-center cursor-pointer transition-all duration-300 ${
                isDragActive ? 'border-cyber-blue bg-cyber-blue/10' : 'hover:border-cyber-blue/50'
              }`}
            >
              <input {...getInputProps()} />
              <Upload className="w-12 h-12 text-cyber-blue mx-auto mb-4" />
              <p className="text-lg mb-2">
                {isDragActive ? 'Drop your CSV file here' : 'Drag & drop CSV file or click to browse'}
              </p>
              <p className="text-sm text-gray-400">Supports .csv and .txt files</p>
            </div>

            {/* Options Panel */}
            {showOptions && (
              <div className="cyber-border p-6 space-y-4">
                <h4 className="font-semibold text-cyber-purple">Conversion Options</h4>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium mb-1">Delimiter</label>
                    <select
                      value={options.delimiter}
                      onChange={(e) => setOptions({...options, delimiter: e.target.value})}
                      className="cyber-input w-full"
                    >
                      <option value=",">Comma (,)</option>
                      <option value=";">Semicolon (;)</option>
                      <option value="\t">Tab</option>
                      <option value="|">Pipe (|)</option>
                    </select>
                  </div>
                  <div>
                    <label className="block text-sm font-medium mb-1">Output Format</label>
                    <select
                      value={options.output_format}
                      onChange={(e) => setOptions({...options, output_format: e.target.value})}
                      className="cyber-input w-full"
                    >
                      <option value="array">Array of Objects</option>
                      <option value="object">Object with Arrays</option>
                    </select>
                  </div>
                </div>
                <div className="flex flex-wrap gap-4">
                  <label className="flex items-center space-x-2">
                    <input
                      type="checkbox"
                      checked={options.has_header}
                      onChange={(e) => setOptions({...options, has_header: e.target.checked})}
                      className="rounded border-cyber-blue/30"
                    />
                    <span className="text-sm">Has Header Row</span>
                  </label>
                  <label className="flex items-center space-x-2">
                    <input
                      type="checkbox"
                      checked={options.pretty_print}
                      onChange={(e) => setOptions({...options, pretty_print: e.target.checked})}
                      className="rounded border-cyber-blue/30"
                    />
                    <span className="text-sm">Pretty Print</span>
                  </label>
                  <label className="flex items-center space-x-2">
                    <input
                      type="checkbox"
                      checked={options.infer_types}
                      onChange={(e) => setOptions({...options, infer_types: e.target.checked})}
                      className="rounded border-cyber-blue/30"
                    />
                    <span className="text-sm">Infer Types</span>
                  </label>
                </div>
              </div>
            )}

            {/* CSV Text Area */}
            <div>
              <textarea
                value={csvData}
                onChange={(e) => setCsvData(e.target.value)}
                placeholder="Paste your CSV data here or drop a file above..."
                className="cyber-input w-full h-64 resize-none font-mono text-sm"
              />
            </div>

            {/* Convert Button */}
            <button
              onClick={convertCSV}
              disabled={!csvData.trim() || loading}
              className="cyber-button w-full py-3 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center space-x-2"
            >
              {loading ? (
                <>
                  <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-white"></div>
                  <span>Converting...</span>
                </>
              ) : (
                <>
                  <Zap className="w-5 h-5" />
                  <span>Convert to JSON</span>
                </>
              )}
            </button>
          </div>

          {/* Output Section */}
          <div className="space-y-6">
            <div className="flex items-center justify-between">
              <h3 className="text-xl font-semibold text-cyber-purple glow-text">JSON Output</h3>
              {jsonOutput && (
                <button
                  onClick={downloadJSON}
                  className="flex items-center space-x-2 px-4 py-2 bg-dark-card border border-cyber-purple/30 rounded hover:border-cyber-purple transition-colors"
                >
                  <Download className="w-4 h-4" />
                  <span>Download</span>
                </button>
              )}
            </div>

            {/* JSON Output Area */}
            <div className="cyber-border h-96 overflow-auto">
              <pre className="p-4 text-sm font-mono whitespace-pre-wrap">
                {jsonOutput || (
                  <span className="text-gray-500 italic">
                    JSON output will appear here after conversion...
                  </span>
                )}
              </pre>
            </div>

            {/* CLI Instructions */}
            <div className="cyber-border p-6">
              <h4 className="font-semibold text-cyber-blue mb-3 flex items-center">
                <Terminal className="w-5 h-5 mr-2" />
                CLI Usage
              </h4>
              <div className="bg-dark-bg rounded p-3 font-mono text-sm">
                <div className="text-cyber-blue"># Install and use CLI</div>
                <div className="text-gray-300">go build -o csv2json</div>
                <div className="text-gray-300">./csv2json -i input.csv -o output.json</div>
                <div className="mt-2 text-cyber-blue"># Start API server</div>
                <div className="text-gray-300">./csv2json -server</div>
              </div>
            </div>
          </div>
        </div>

        {/* Footer */}
        <footer className="mt-16 text-center text-gray-400 border-t border-dark-border pt-8">
          <p>Built with Go, React, and Tailwind CSS • Open Source • Lightning Fast</p>
        </footer>
      </div>
    </div>
  );
}

export default App;
