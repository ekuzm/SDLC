"""Load and validate SysML sources in Pilot, then export the eight model views.

Layout follows the approach of kiryshabutor/SDLC#1: %viz emits PlantUML;
invisible edges arrange Row children without changing model semantics.
No HTML rendering or hand-drawn diagram geometry is used.
"""
from pathlib import Path
import hashlib
import json
import re
import subprocess
import cairosvg
import nbformat
from nbclient import NotebookClient

ROOT = Path(__file__).resolve().parents[1]
OUT = ROOT.parent / 'images'
VIEWS = json.loads((ROOT / 'tools/views.json').read_text())

def arrange_rows(source):
    stack = ['root']
    children = {'root': []}
    horizontal = set()
    for line in source.splitlines():
        match = re.match(r'^rec\s+\w+\s+"(.*)"\s+as\s+(E\d+)', line)
        if match:
            caption, ident = match.groups()
            children[stack[-1]].append(ident)
            children[ident] = []
            if line.rstrip().endswith('{'):
                stack.append(ident)
                if re.search(r': (Row|NavBar|ListItem)$', caption):
                    horizontal.add(ident)
        elif line.strip() == '}' and len(stack) > 1:
            stack.pop()

    def leaf(ident, last=False):
        while children[ident]:
            ident = children[ident][-1 if last else 0]
        return ident

    edges = []
    for ident, elements in children.items():
        direction = 'right' if ident in horizontal else 'down'
        for left, right in zip(elements, elements[1:]):
            edges.append(f'{leaf(left, ident not in horizontal)} -[hidden]{direction}-> {leaf(right)}')
    source = re.sub(r'"<\w+> ', '"', source)
    return source.replace('@enduml', '\n'.join(edges) + '\n@enduml')

def output_text(cell):
    return '\n'.join(str(o.get('text', o.get('data', {}).get('text/plain', ''))) for o in cell.outputs)

def digest(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()

def main():
    files = [ROOT/'balance.sysml', ROOT/'ui-kit.sysml'] + [ROOT/v['file'] for v in VIEWS]
    notebook = nbformat.v4.new_notebook()
    # One compilation unit permits forward references across model files.
    notebook.cells.append(nbformat.v4.new_code_cell('\n\n'.join(p.read_text() for p in files)))
    for view in VIEWS:
        style = '--style=PUMLCODE'
        notebook.cells.append(nbformat.v4.new_code_cell(f"%viz --view={view['view']} {style} {view['element']}"))
    NotebookClient(notebook, kernel_name='sysml', timeout=600, allow_errors=False).execute()
    text = output_text(notebook.cells[0])
    if re.search(r'\bERROR\b|\bError\b', text) or any(o.output_type == 'error' for o in notebook.cells[0].outputs):
        raise RuntimeError('SysML validation failed:\n' + text)
    print('Models loaded:', text.strip(), flush=True)
    OUT.mkdir(exist_ok=True)
    for view, cell in zip(VIEWS, notebook.cells[1:]):
        data = {k: v for o in cell.outputs for k, v in o.get('data', {}).items()}
        code = data.get('text/plain', '')
        if not code.startswith('@startuml'):
            raise RuntimeError(f"No PlantUML for {view['name']}: {output_text(cell)}")
        if view['wire']:
            code = arrange_rows(code)
        elif view['view'] == 'case':
            code = code.replace('@startuml', '@startuml\nleft to right direction', 1)
        code = re.sub(r'"<\w+> ', '"', code)
        # Explicit font keeps Java layout and Cairo rasterization metrics consistent.
        code = code.replace('@enduml', 'skinparam defaultFontName DejaVu Sans\nskinparam padding 5\nhide empty members\n@enduml')
        result = subprocess.run(['java', '-cp', '/opt/conda/share/jupyter/kernels/sysml/*',
                                 'net.sourceforge.plantuml.Run', '-tsvg', '-pipe', '-charset', 'UTF-8'],
                                input=code.encode(), capture_output=True, check=True)
        svg = result.stdout
        if not svg or b'<svg' not in svg:
            raise RuntimeError(f"No SVG for {view['name']}: {output_text(cell)}")
        cairosvg.svg2png(bytestring=svg, write_to=str(OUT/(view['name']+'.png')),
                        scale=1.5, background_color='white')
        print(view['name'] + ': PNG exported', flush=True)
    report = {'tool': 'SysML v2 Pilot Implementation 0.62.0',
              'validation': 'All source files loaded together without kernel errors',
              'sources': {str(p.relative_to(ROOT)): digest(p) for p in files},
              'tooling': {p.name: digest(p) for p in (ROOT/'tools').iterdir() if p.name in ['render.py', 'render-all.sh', 'Dockerfile', 'views.json']},
              'images': {v['name']+'.png': digest(OUT/(v['name']+'.png')) for v in VIEWS}}
    (ROOT/'tools/render-manifest.json').write_text(json.dumps(report, ensure_ascii=False, indent=2)+'\n')

if __name__ == '__main__':
    main()
