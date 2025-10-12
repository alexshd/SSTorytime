import './style.css';

// Load main layout from app.html and initialize app
fetch('/src/app.html')
  .then(res => res.text())
  .then(html =>
  {
    document.querySelector('#app').innerHTML = html;

    // Select DOM elements
    const inputText = document.querySelector('#input-text');
    const outputArea = document.querySelector('#output-area');
    const convertBtn = document.querySelector('#convert-btn');
    const copyBtn = document.querySelector('#copy-btn');
    const uploadBtn = document.querySelector('#upload-btn');
    const fileInput = document.querySelector('#file-input');
    const arrowsBtn = document.querySelector('#arrows-btn');
    const saveBtn = document.querySelector('#save-btn');
    const statusMessage = document.querySelector('#status-message');
    const titleSection = document.querySelector('#title-section');
    const arrowsSection = document.querySelector('#arrows-section');
    const arrowsContent = document.querySelector('#arrows-content');

    let currentFileName = '';
    let isScrollSyncing = false;

    // === ALL FUNCTION DEFINITIONS MUST BE HERE ===

    function showStatus(message, isError = false)
    {
      const p = statusMessage.querySelector('p');
      p.textContent = message;
      p.className = 'text-sm ' + (isError ? 'text-red-600' : 'text-green-600');
      statusMessage.classList.remove('hidden');
      setTimeout(() => { statusMessage.classList.add('hidden'); }, 3000);
    }

    async function convertText()
    {
      const text = inputText.value.trim();

      if (!text)
      {
        showStatus('Please enter some text to convert', true);
        return;
      }

      convertBtn.disabled = true;
      convertBtn.textContent = 'Converting...';

      try
      {
        const formData = new FormData();
        formData.append('text', text);

        const response = await fetch('/api/convert', {
          method: 'POST',
          body: formData
        });

        if (!response.ok)
        {
          throw new Error('API Error: ' + response.status);
        }

        const n4lOutput = await response.text();
        outputArea.innerHTML = highlightArrows(n4lOutput);
        copyBtn.disabled = false;
        saveBtn.disabled = false;
        arrowsBtn.disabled = false;
        showStatus('Text converted successfully!');

      } catch (error)
      {
        console.error('Conversion error:', error);
        showStatus('Error: ' + error.message, true);
        outputArea.innerHTML = '';
        copyBtn.disabled = true;
        saveBtn.disabled = true;
        arrowsBtn.disabled = true;
      } finally
      {
        convertBtn.disabled = false;
        convertBtn.textContent = 'Convert to N4L';
      }
    }

    async function copyToClipboard()
    {
      try
      {
        await navigator.clipboard.writeText(outputArea.innerText);
        showStatus('Copied to clipboard!');
      } catch (error)
      {
        console.error('Copy error:', error);
        showStatus('Failed to copy to clipboard', true);
      }
    }

    function handleFileUpload()
    {
      fileInput.click();
    }

    function handleFileSelect(event)
    {
      const file = event.target.files[0];

      if (!file)
      {
        return;
      }

      // Check file size (limit to 10MB)
      if (file.size > 10 * 1024 * 1024)
      {
        showStatus('File too large. Please select a file smaller than 10MB.', true);
        return;
      }

      // We'll attempt to read any file as text
      // The browser will handle it gracefully if it's not a text file
      const reader = new FileReader();

      reader.onload = function (e)
      {
        const content = e.target.result;
        inputText.value = content;
        currentFileName = file.name;

        // Hide title section when file is uploaded
        titleSection.classList.add('hidden');

        showStatus('File "' + file.name + '" loaded successfully!');

        // Clear any previous output
        outputArea.innerText = '';
        copyBtn.disabled = true;
        saveBtn.disabled = true;
        arrowsBtn.disabled = true;
        arrowsSection.classList.add('hidden');

        // Reset file input
        fileInput.value = '';
      }; reader.onerror = function ()
      {
        showStatus('Error reading file. Please try again.', true);
        fileInput.value = '';
      };

      reader.readAsText(file);
    }

    // Complete list of valid N4L arrows from SSTconfig files
    // This is used for validation - if an arrow is not in this list, it's invalid
    function getValidArrowsList()
    {
      return [
        // NR-0: Similarity arrows (symmetric)
        'similar to', 'sim', 'associated with', 'ass', 'see also', 'see', 'alias',
        'equals', '=', 'compare to', 'compare', 'is not', 'not the same as',
        'looks like', 'll', 'sounds like', 'sl', 'don\'t confuse with', 'confuse',
        'alternate spelling', 'sp', 'stands for', 'sfor', 'is sometimes mistaken for', 'mistaken',
        'has same supplier', 'sm-sup', 'has common origin', 'cm-org', 'is entangled with', 'entg',
        'has same destination', 'sm-dir', 'same class', 'sm_cls', 'same group', 'sm_grp',
        'same tribe', 'sm_trb', 'same hometown', 'sm_hmt', 'same modus operandi', 'sm_MO',
        'is a variant of', 'variant', 'near to', 'nr', 'terminates together with', 'termwith',

        // LT-1: Causality arrows (directional)
        'leads to', 'fwd', 'comes from', 'bwd', 'brings about', 'brings', 'was brought about by', 'brought-by',
        'from which we derive', 'derive', 'derives from', 'derive-fr', 'leads ahead to', 'fwd1',
        'leads back to', 'bwd1', 'succeeds', 'succ1', 'is succeeded by', 'succ-by',
        'precedes', 'prec', 'is preceded by', 'prec-by', 'succeeded by', 'succ2',
        'preceded by', 'pre-by', 'comes before', 'bfr', 'comes after', 'after',
        'supplies', 'suppl', 'has been supplied by', 'suppl-by', 'evolved to', 'evolve',
        'evolved from', 'evolve-fr', 'thereafter', '=>', 'whence', '<=',
        'next if yes', 'ifyes', 'is affirmative outcome of', 'bifyes',
        'next if no', 'if no', 'is negitive outcome of', 'bifno',
        'affects', 'aff', 'is affected by', 'aff-by', 'causes', 'cause',
        'is caused by', 'cause-by', 'creates', 'cr', 'is created by', 'crtd',
        'redirects', 'redir', 'is redirected by', 'redir-by',
        'is used by', 'used-by', 'makes use of', 'uses', 'binds to', 'bind-to',
        'is bound by', 'bound-by', 'results in', 'result', 'was a result of', 'result-of',
        'enables', 'depends on', 'dep', 'is depended upon by', 'dep-by',
        'invokes', 'invoke', 'is invoked by', 'invoke-by', 'determines', 'det',

        // CN-2: Composition arrows
        'contains', 'contain', 'belongs to', 'belong', 'consists of', 'consists',
        'makes up part of', 'mkpt', 'uses word', 'useword', 'is a word used in', 'word-in',
        'uses material', 'material', 'is a material for', 'mat-for',
        'may contain', 'm-cont', 'may be contained by', 'm-cont-by',
        'has component', 'has-cmpt', 'is component of', 'cmpt-of',
        'has a part', 'has-pt', 'is a part of', 'pt-of',
        'has no part', 'no-pt', 'is not part of', 'not-pt-of',
        'has ingredient', 'ingred', 'is an ingredient of', 'ingred-of',
        'does not contain', 'has no', 'is not from set', 'is not',
        'is a set of', 'setof', 'is part of the set', 'in-set',
        'is a group of', 'grp-of', 'is an element of', 'in-grp',
        'is a herd of', 'herd', 'is part of herd', 'in-herd',
        'is a tribe of', 'tribe', 'is part of tribe', 'in-tribe',
        'is a region of', 'region', 'is within region', 'within',
        'is an interval of', 'interval', 'is part of interval', 'in-intvl',
        'species has member', 'smember', 'belongs to species', 'sbelongs',
        'is an umbrella for', 'umbrella', 'under the aegis of', 'aegis',
        'is the company of', 'company', 'works for', 'works',
        'is the employer of', 'employer', 'is an employee of', 'employee',
        'has many copies of', 'has-many', 'is one of many in', 'one-of',
        'has member', 'has-memb', 'is a member of', 'is-memb',
        'subsumes', 'sub', 'is subsumed by', 'sub-by', 'encloses', 'encl',
        'is enclosed by', 'is-encl', 'has absorbed', 'sw', 'absorbed by', 'ab-by',
        'houses', 'house', 'is housed by', 'is-housed', 'situates', 'situate',
        'is situated in', 'is-sit', 'bullet point', 'bullet', 'is a bulletpoint for', 'isbullet',
        'occurs within', 'occurs', 'situates occurrences of', 'occurence-of',
        'is a shared resource for', 'shared-by', 'shares resource', 'shares',

        // EP-3: Properties/Expresses arrows
        'note', 'remark', 'is a note or remark about', 'isnotefor',
        'added remark', 'is a remark about', 'remabout', 'please note!', 'NB',
        'is an important remark concerning', 'important regarding',
        'has example', 'e.g.', 'is an example of', 'ex', 'has no example', 'noex',
        'is not an example of', 'not-ex', 'has abbreviation', 'abbrev',
        'is short for', 'short for', 'has version', 'version', 'is a version of', 'is-version',
        'has supplier', 'supply', 'is the supplier of', 'supply-by',
        'has source', 'source', 'the source for', 'source-of',
        'discusses', 'disc', 'is discussed in', 'is-disc',
        'described as', 'descr', 'is used as a description of', 'isdescr',
        'has age character', 'age', 'is the age characteristic of', 'isageof',
        'mentions study', 'ment', 'a kind of study mentioned in', 'isment',
        'has a pun about', 'pun on', 'was made into pun', 'has pun',
        'famously uses', 'fuses', 'famously used in', 'fused',
        'has title', 'title', 'is the title of', 'title of',
        'has description', 'description', 'is the description for', 'descrip-for',
        'is about topic/them', 'about', 'is the topic/theme of', 'theme-of',
        'is wrong about', 'wrongabt', 'is not accurately represented by', 'narb',
        'expresses', 'expresses property', 'expr', 'is/are expressed by', 'expr-by',
        'has value', 'hasX', 'is the value of', 'isXof',
        'may have value', 'maybeX', 'is the possible value for', 'mXof',
        'has frequency', 'freq', 'is the frequency of', 'freq-of',
        'is characterized by', 'state', 'is a/the state of', 'state-of',
        'is expressed in formula', 'formula', 'is a formula for', 'isformula',
        'has condition', 'cond', 'is the condition of', 'cond-of',
        'has an attribute', 'attr', 'is an attribute of', 'attr-of',
        'has feature', 'feat', 'is a feature of', 'feat-of',
        'has property', 'prop', 'is a property of', 'prop-of',

        // Special arrows
        'offers', 'accepts', 'observes', 'ovserves', 'may be influenced by',
        'has overlap', 'has promiser', 'has promisee', 'observed by',
        'has intended promisee', 'has intended promiser', 'imposes'
      ];
    }

    function isValidArrow(arrowText)
    {
      // Remove parentheses and normalize
      const normalized = arrowText.replace(/^\(|\)$/g, '').trim().toLowerCase();
      const validArrows = getValidArrowsList();

      // Check if it's a valid arrow (exact match or close match)
      return validArrows.some(function (valid)
      {
        return valid.toLowerCase() === normalized;
      });
    }

    function highlightArrows(text)
    {
      // First escape HTML entities to preserve text as-is
      let escaped = text
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#39;');

      // Only highlight parenthetical arrow descriptions like "(appears close to)"
      // This matches any text in parentheses that looks like a relationship description
      const parenArrowRegex = /\([a-z][a-z\s,;:.\-'\/]*\)/gi;

      // Highlight parenthetical arrow descriptions
      const result = escaped.replace(parenArrowRegex, function (match)
      {
        // Only highlight if it's not just a single word (to avoid matching things like (x) or (a))
        // and if it contains at least one space or common relationship words
        if (match.length > 4 && (match.includes(' ') || /(?:to|by|from|in|on|at|as|with|of)/i.test(match)))
        {
          // Check if this is a valid N4L arrow
          const isValid = isValidArrow(match);
          const cssClass = isValid ? 'n4l-arrow-highlight' : 'n4l-arrow-error';

          return '<span class="' + cssClass + '" data-arrow="' + match + '" data-valid="' + isValid + '" onclick="showArrowMenu(event, this)">' +
            match +
            '</span>';
        }
        return match; // Don't highlight short or non-relational parenthetical text
      });

      return result;
    }

    // Global function for arrow menu (called from onclick in HTML)
    window.showArrowMenu = function (event, arrowSpan)
    {
      event.stopPropagation();

      // Remove any existing menu
      const existingMenu = document.querySelector('.arrow-menu');
      if (existingMenu) existingMenu.remove();

      const arrow = arrowSpan.dataset.arrow;
      const isParenthetical = arrow.startsWith('(');

      // Create menu
      const menu = document.createElement('div');
      menu.className = 'arrow-menu';
      menu.style.cssText = 'position: absolute; z-index: 100; background: white; border: 1px solid #38bdf8; border-radius: 0.5rem; box-shadow: 0 4px 12px rgba(0,0,0,0.15); padding: 0.75rem; min-width: 280px; max-height: 400px; overflow-y: auto;';

      // N4L Semantic Arrow Types based on SSTconfig
      const allArrows = [
        // NR-0: Similarity (symmetric, non-directional)
        { type: 'NR-0', symbol: '(similar to)', abbrev: '(sim)', keywords: ['similar', 'like', 'resembles', 'appears close to'] },
        { type: 'NR-0', symbol: '(associated with)', abbrev: '(ass)', keywords: ['associated', 'related', 'connected'] },
        { type: 'NR-0', symbol: '(see also)', abbrev: '(see)', keywords: ['see', 'refer', 'also', 'related'] },
        { type: 'NR-0', symbol: '(alias)', abbrev: '(alias)', keywords: ['also called', 'known as', 'alias'] },
        { type: 'NR-0', symbol: '(equals)', abbrev: '(=)', keywords: ['equals', 'same as', 'equivalent', 'is'] },
        { type: 'NR-0', symbol: '(compare to)', abbrev: '(compare)', keywords: ['compare', 'contrast', 'versus'] },
        { type: 'NR-0', symbol: '(is not)', abbrev: '(!eq)', keywords: ['not', 'different', 'unlike', 'isn\'t'] },

        // LT-1: Leads To (causality, processes)
        { type: 'LT-1', symbol: '(leads to)', abbrev: '(fwd)', keywords: ['leads', 'causes', 'results', 'produces', 'brings about'] },
        { type: 'LT-1', symbol: '(causes)', abbrev: '(cause)', keywords: ['causes', 'triggers', 'initiates', 'starts'] },
        { type: 'LT-1', symbol: '(creates)', abbrev: '(cr)', keywords: ['creates', 'generates', 'makes', 'builds'] },
        { type: 'LT-1', symbol: '(results in)', abbrev: '(result)', keywords: ['results', 'ends', 'concludes', 'becomes'] },
        { type: 'LT-1', symbol: '(enables)', abbrev: '(enables)', keywords: ['enables', 'allows', 'permits', 'facilitates'] },
        { type: 'LT-1', symbol: '(affects)', abbrev: '(aff)', keywords: ['affects', 'impacts', 'influences', 'changes'] },
        { type: 'LT-1', symbol: '(precedes)', abbrev: '(prec)', keywords: ['precedes', 'before', 'prior', 'earlier'] },
        { type: 'LT-1', symbol: '(comes from)', abbrev: '(bwd)', keywords: ['from', 'derives', 'originates', 'stems'] },
        { type: 'LT-1', symbol: '(determines)', abbrev: '(det)', keywords: ['determines', 'decides', 'controls', 'defines'] },

        // CN-2: Contains (composition, part-whole)
        { type: 'CN-2', symbol: '(contains)', abbrev: '(contain)', keywords: ['contains', 'includes', 'has', 'comprises'] },
        { type: 'CN-2', symbol: '(is a part of)', abbrev: '(pt-of)', keywords: ['part of', 'within', 'inside', 'component of'] },
        { type: 'CN-2', symbol: '(consists of)', abbrev: '(consists)', keywords: ['consists', 'made up', 'composed'] },
        { type: 'CN-2', symbol: '(belongs to)', abbrev: '(belong)', keywords: ['belongs', 'member', 'element'] },
        { type: 'CN-2', symbol: '(has component)', abbrev: '(has-cmpt)', keywords: ['component', 'piece', 'section'] },
        { type: 'CN-2', symbol: '(is a set of)', abbrev: '(setof)', keywords: ['set', 'collection', 'group'] },

        // EP-3: Expresses/Properties (attributes, descriptions)
        { type: 'EP-3', symbol: '(expresses)', abbrev: '(expr)', keywords: ['expresses', 'shows', 'demonstrates', 'reveals'] },
        { type: 'EP-3', symbol: '(note)', abbrev: '(note)', keywords: ['note', 'remark', 'comment', 'annotation'] },
        { type: 'EP-3', symbol: '(defined as)', abbrev: '(def)', keywords: ['defined', 'means', 'definition', 'is defined'] },
        { type: 'EP-3', symbol: '(has example)', abbrev: '(e.g.)', keywords: ['example', 'instance', 'such as', 'for example'] },
        { type: 'EP-3', symbol: '(is about)', abbrev: '(about)', keywords: ['about', 'concerning', 'regarding', 'topic'] },
        { type: 'EP-3', symbol: '(has attribute)', abbrev: '(attr)', keywords: ['attribute', 'property', 'characteristic', 'feature'] },
        { type: 'EP-3', symbol: '(has feature)', abbrev: '(feat)', keywords: ['feature', 'aspect', 'quality'] },
        { type: 'EP-3', symbol: '(describes)', abbrev: '(descr)', keywords: ['describes', 'explains', 'details'] },

        // Special annotations
        { type: 'Special', symbol: '(is a special case of)', abbrev: '(**)', keywords: ['special case', 'exception', 'specific'] },
        { type: 'Special', symbol: '(is an example of)', abbrev: '(>>)', keywords: ['example of', 'instance of', 'exemplifies'] },
        { type: 'Special', symbol: '(discusses)', abbrev: '(%)', keywords: ['discusses', 'talks about', 'addresses'] },
        { type: 'Special', symbol: '(involves)', abbrev: '(=)', keywords: ['involves', 'includes', 'entails'] }
      ];

      let menuHTML = '';

      // Check if this arrow is invalid and show error message
      const isValid = arrowSpan.dataset.valid === 'true';
      if (!isValid && isParenthetical)
      {
        menuHTML += '<div style="background: #fee2e2; color: #991b1b; padding: 0.75rem; border-radius: 0.375rem; margin-bottom: 0.75rem; border-left: 4px solid #dc2626;">' +
          '<div style="font-weight: 700; font-size: 0.875rem; margin-bottom: 0.25rem;">⚠️ Invalid N4L Arrow</div>' +
          '<div style="font-size: 0.75rem; line-height: 1.4;">This arrow is not recognized in the N4L configuration. ' +
          'Please select a valid arrow from the list below or the N4L parser will report an error.</div>' +
          '</div>';
      }

      if (isParenthetical)
      {
        // For parenthetical text, suggest matching arrows
        const lowerArrow = arrow.toLowerCase();
        const suggestedArrows = allArrows.filter(function (arr)
        {
          return arr.keywords.some(function (keyword)
          {
            return lowerArrow.includes(keyword);
          });
        });

        menuHTML += '<div style="font-weight: 700; color: #0369a1; margin-bottom: 0.75rem; font-size: 0.875rem; border-bottom: 2px solid #38bdf8; padding-bottom: 0.5rem;">Replace with N4L Arrow:</div>';

        // Show suggested arrows first
        if (suggestedArrows.length > 0)
        {
          menuHTML += '<div style="color: #16a34a; font-size: 0.75rem; margin-bottom: 0.5rem; font-weight: 600; text-transform: uppercase;">✓ Suggested Matches:</div>';
          suggestedArrows.forEach(function (arr)
          {
            menuHTML += '<div class="arrow-menu-item" style="padding: 0.5rem; cursor: pointer; border-radius: 0.375rem; font-size: 0.875rem; margin-bottom: 0.25rem; background: #f0fdf4; border-left: 3px solid #22c55e;" ' +
              'onmouseover="this.style.background=\'#dcfce7\'" ' +
              'onmouseout="this.style.background=\'#f0fdf4\'" ' +
              'onclick="replaceArrow(\'' + arrow.replace(/'/g, "\\'") + '\', \'' + arr.symbol.replace(/'/g, "\\'") + '\', this.closest(\'.arrow-menu\'))">' +
              '<div style="display: flex; justify-content: space-between; align-items: center;">' +
              '<span style="font-family: monospace; color: #0369a1; font-weight: 600; flex: 1;">' + arr.symbol + '</span>' +
              '<span style="font-size: 0.75rem; color: #64748b; margin-left: 0.5rem;">' + arr.abbrev + '</span>' +
              '</div>' +
              '<div style="font-size: 0.7rem; color: #059669; margin-top: 0.25rem;">Type: ' + arr.type + '</div>' +
              '</div>';
          });

          menuHTML += '<div style="color: #64748b; font-size: 0.75rem; margin: 0.75rem 0 0.5rem; font-weight: 600; text-transform: uppercase;">Browse All Arrows:</div>';
        }

        // Group other arrows by type
        const arrowsByType = {};
        const otherArrows = suggestedArrows.length > 0
          ? allArrows.filter(function (arr) { return !suggestedArrows.includes(arr); })
          : allArrows;

        otherArrows.forEach(function (arr)
        {
          if (!arrowsByType[arr.type]) arrowsByType[arr.type] = [];
          arrowsByType[arr.type].push(arr);
        });

        // Display arrows grouped by type
        const typeOrder = ['NR-0', 'LT-1', 'CN-2', 'EP-3', 'Special'];
        const typeNames = {
          'NR-0': '🔗 Similarity (Non-directional)',
          'LT-1': '➡️ Causality (Leads To)',
          'CN-2': '📦 Composition (Contains)',
          'EP-3': '🏷️ Properties (Expresses)',
          'Special': '⭐ Special Annotations'
        };

        typeOrder.forEach(function (type)
        {
          if (arrowsByType[type] && arrowsByType[type].length > 0)
          {
            menuHTML += '<div style="font-size: 0.7rem; color: #475569; margin: 0.75rem 0 0.25rem; font-weight: 600;">' + typeNames[type] + '</div>';
            arrowsByType[type].forEach(function (arr)
            {
              menuHTML += '<div class="arrow-menu-item" style="padding: 0.4rem 0.5rem; cursor: pointer; border-radius: 0.25rem; font-size: 0.8rem; margin-bottom: 0.15rem;" ' +
                'onmouseover="this.style.background=\'#e0f2fe\'" ' +
                'onmouseout="this.style.background=\'white\'" ' +
                'onclick="replaceArrow(\'' + arrow.replace(/'/g, "\\'") + '\', \'' + arr.symbol.replace(/'/g, "\\'") + '\', this.closest(\'.arrow-menu\'))">' +
                '<div style="display: flex; justify-content: space-between; align-items: center;">' +
                '<span style="font-family: monospace; color: #0369a1; font-weight: 500; font-size: 0.75rem;">' + arr.symbol + '</span>' +
                '<span style="font-size: 0.65rem; color: #94a3b8;">' + arr.abbrev + '</span>' +
                '</div>' +
                '</div>';
            });
          }
        });
      }
      else
      {
        // For arrow symbols, show options to change to other arrows grouped by type
        menuHTML += '<div style="font-weight: 700; color: #0369a1; margin-bottom: 0.75rem; font-size: 0.875rem; border-bottom: 2px solid #38bdf8; padding-bottom: 0.5rem;">Change N4L Arrow:</div>';

        // Group arrows by type
        const arrowsByType = {};
        allArrows.forEach(function (arr)
        {
          if (arr.symbol !== arrow)
          {
            if (!arrowsByType[arr.type]) arrowsByType[arr.type] = [];
            arrowsByType[arr.type].push(arr);
          }
        });

        const typeOrder = ['NR-0', 'LT-1', 'CN-2', 'EP-3', 'Special'];
        const typeNames = {
          'NR-0': '🔗 Similarity (NR-0)',
          'LT-1': '➡️ Causality (LT-1)',
          'CN-2': '📦 Composition (CN-2)',
          'EP-3': '🏷️ Properties (EP-3)',
          'Special': '⭐ Special'
        };

        typeOrder.forEach(function (type)
        {
          if (arrowsByType[type] && arrowsByType[type].length > 0)
          {
            menuHTML += '<div style="font-size: 0.7rem; color: #475569; margin: 0.75rem 0 0.25rem; font-weight: 600;">' + typeNames[type] + '</div>';
            arrowsByType[type].forEach(function (arr)
            {
              menuHTML += '<div class="arrow-menu-item" style="padding: 0.4rem 0.5rem; cursor: pointer; border-radius: 0.25rem; font-size: 0.8rem; margin-bottom: 0.15rem;" ' +
                'onmouseover="this.style.background=\'#e0f2fe\'" ' +
                'onmouseout="this.style.background=\'white\'" ' +
                'onclick="replaceArrow(\'' + arrow.replace(/'/g, "\\'") + '\', \'' + arr.symbol.replace(/'/g, "\\'") + '\', this.closest(\'.arrow-menu\'))">' +
                '<div style="display: flex; justify-content: space-between; align-items: center;">' +
                '<span style="font-family: monospace; color: #0369a1; font-weight: 500; font-size: 0.75rem;">' + arr.symbol + '</span>' +
                '<span style="font-size: 0.65rem; color: #94a3b8;">' + arr.abbrev + '</span>' +
                '</div>' +
                '</div>';
            });
          }
        });
      }

      menuHTML += '<hr style="margin: 0.75rem 0; border-color: #e2e8f0;">';
      menuHTML += '<div class="arrow-menu-item" style="padding: 0.5rem; cursor: pointer; border-radius: 0.375rem; color: #dc2626; font-size: 0.875rem; font-weight: 600; text-align: center;" ' +
        'onmouseover="this.style.background=\'#fee2e2\'" ' +
        'onmouseout="this.style.background=\'white\'" ' +
        'onclick="deleteArrow(\'' + arrowSpan.dataset.arrow.replace(/'/g, "\\'") + '\', this.closest(\'.arrow-menu\'))">🗑️ Delete entire line</div>';

      menu.innerHTML = menuHTML;


      // Position menu near the arrow
      const rect = arrowSpan.getBoundingClientRect();
      menu.style.left = rect.left + 'px';
      menu.style.top = (rect.bottom + 5) + 'px';

      document.body.appendChild(menu);

      // Close menu when clicking outside
      setTimeout(function ()
      {
        document.addEventListener('click', function closeMenu()
        {
          const m = document.querySelector('.arrow-menu');
          if (m) m.remove();
          document.removeEventListener('click', closeMenu);
        });
      }, 100);
    };

    // Replace arrow in the output
    window.replaceArrow = function (oldArrow, newArrow, menu)
    {
      const text = outputArea.innerHTML;

      // Escape special characters for regex
      function escapeRegex(str)
      {
        return str.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
      }

      // Escape HTML entities in the old arrow
      const oldEscaped = oldArrow
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#39;');

      // Escape HTML entities in the new arrow
      const newEscaped = newArrow
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#39;');

      // Find and replace the span containing the old arrow
      const regex = new RegExp('<span class="n4l-arrow-highlight"[^>]*data-arrow="[^"]*"[^>]*>' + escapeRegex(oldEscaped) + '</span>');
      const newText = text.replace(regex, '<span class="n4l-arrow-highlight" data-arrow="' + newArrow + '" onclick="showArrowMenu(event, this)">' + newEscaped + '</span>');
      outputArea.innerHTML = newText;

      menu.remove();
      showStatus('Arrow changed to ' + newArrow);
    };

    // Delete arrow from the output
    window.deleteArrow = function (arrow, menu)
    {
      const text = outputArea.innerHTML;

      // Escape special characters for regex
      function escapeRegex(str)
      {
        return str.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
      }

      // Escape HTML entities
      const escaped = arrow
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#39;');

      // Find the line containing the arrow and delete the entire line
      // Match from start of line (or after a <br> or newline) to end of line (or before <br> or newline)
      const lineRegex = new RegExp(
        '(?:^|(?:<br>)|(?:\\n))([^\\n<]*<span class="n4l-arrow-highlight"[^>]*data-arrow="[^"]*"[^>]*>' +
        escapeRegex(escaped) +
        '</span>[^\\n<]*)(?:(?:<br>)|(?:\\n)|$)',
        'g'
      );

      let newText = text.replace(lineRegex, '');

      // If that didn't work (no <br> tags), try a simpler approach - just remove the whole line
      if (newText === text)
      {
        // Split by line breaks, filter out lines containing the arrow, rejoin
        const lines = text.split(/\n/);
        const filteredLines = lines.filter(function (line)
        {
          return !line.includes('<span class="n4l-arrow-highlight"[^>]*data-arrow="[^"]*"[^>]*>' + escapeRegex(escaped) + '</span>');
        });
        newText = filteredLines.join('\n');
      }

      outputArea.innerHTML = newText;

      menu.remove();
      showStatus('Line with arrow deleted');
    };

    function extractArrows()
    {
      // Toggle arrows section visibility
      const isHidden = arrowsSection.classList.contains('hidden');

      if (isHidden)
      {
        const n4lText = outputArea.innerText.trim();
        if (!n4lText)
        {
          showStatus('No N4L output to show arrows for', true);
          return;
        }

        // Define available arrow types for editing
        const availableArrows = [
          {
            type: 'Implication',
            symbol: '->',
            description: 'Basic implication or causation',
            example: 'A -> B',
            category: 'Basic'
          },
          {
            type: 'Strong Implication',
            symbol: '+>',
            description: 'Strong or positive implication',
            example: 'A +> B',
            category: 'Basic'
          },
          {
            type: 'Weak Implication',
            symbol: '~>',
            description: 'Weak or uncertain implication',
            example: 'A ~> B',
            category: 'Basic'
          },
          {
            type: 'Bidirectional',
            symbol: '<->',
            description: 'Two-way relationship',
            example: 'A <-> B',
            category: 'Bidirectional'
          },
          {
            type: 'Conditional',
            symbol: '?>',
            description: 'Conditional relationship',
            example: 'A ?> B',
            category: 'Conditional'
          },
          {
            type: 'Temporal',
            symbol: '>>',
            description: 'Temporal sequence',
            example: 'A >> B',
            category: 'Temporal'
          },
          {
            type: 'Contextual',
            symbol: '::>',
            description: 'Context-dependent relationship',
            example: 'A ::> B',
            category: 'Contextual'
          },
          {
            type: 'Negation',
            symbol: '!>',
            description: 'Negative implication',
            example: 'A !> B',
            category: 'Negation'
          },
          {
            type: 'Constraint',
            symbol: '|>',
            description: 'Constraint or requirement',
            example: 'A |> B',
            category: 'Constraint'
          },
          {
            type: 'Optional',
            symbol: '~?>',
            description: 'Optional or possible relationship',
            example: 'A ~?> B',
            category: 'Optional'
          }
        ];

        displayAvailableArrows(availableArrows);
        arrowsSection.classList.remove('hidden');
        showStatus('Showing available arrow types');
      }
      else
      {
        arrowsSection.classList.add('hidden');
        showStatus('Arrow types hidden');
      }
    }

    function displayAvailableArrows(arrows)
    {
      const categories = [...new Set(arrows.map(arrow => arrow.category))];

      let content = '<div class="space-y-4">';

      categories.forEach(category =>
      {
        const categoryArrows = arrows.filter(arrow => arrow.category === category);


        content += '<div class="category-section">' +
          '<h4 class="font-semibold text-blue-700 mb-2 text-sm">' + category + ' Arrows</h4>' +
          '<div class="grid grid-cols-1 md:grid-cols-2 gap-2">';

        categoryArrows.forEach(function (arrow)
        {
          content += '<div class="arrow-item p-3 bg-white rounded border border-gray-200 hover:border-blue-300 cursor-pointer transition-colors" ' +
            'onclick="insertArrowTemplate(\'' + arrow.symbol + '\', \' ' + arrow.example + '\')">' +
            '<div class="flex items-center justify-between mb-1">' +
            '<span class="font-mono text-lg text-blue-600">' + arrow.symbol + '</span>' +
            '<span class="text-xs text-gray-500">' + arrow.type + '</span>' +
            '</div>' +
            '<div class="text-sm text-gray-700 mb-1">' + arrow.description + '</div>' +
            '<div class="text-xs font-mono text-gray-500">' + arrow.example + '</div>' +
            '</div>';
        });

        content += '</div></div>';
      });

      content += '</div>';
      arrowsContent.innerHTML = content;
    }

    function insertArrowTemplate(symbol, example)
    {
      const template = '\n' + example + '\n';
      insertAtCaret(outputArea, template);
      outputArea.focus();
      showStatus('Inserted ' + symbol + ' arrow template');
    }

    function insertAtCaret(editableDiv, text)
    {
      const sel = window.getSelection();
      if (!sel.rangeCount) return;
      const range = sel.getRangeAt(0);
      range.deleteContents();
      const textNode = document.createTextNode(text);
      range.insertNode(textNode);
      range.setStartAfter(textNode);
      range.setEndAfter(textNode);
      sel.removeAllRanges();
      sel.addRange(range);
    }

    window.insertArrowTemplate = insertArrowTemplate;

    function saveAsN4L()
    {
      const content = outputArea.innerText.trim();

      if (!content)
      {
        showStatus('No content to save', true);
        return;
      }

      const fileName = currentFileName
        ? currentFileName.replace(/\.[^/.]+$/, '.n4l')
        : 'output.n4l';

      const blob = new Blob([content], { type: 'text/plain' });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = fileName;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);

      showStatus('File saved as ' + fileName);
    }

    function syncScroll()
    {
      if (isScrollSyncing) return;

      isScrollSyncing = true;

      const outputLines = outputArea.innerText.split('\n');
      const inputLines = inputText.value.split('\n');

      // Calculate scroll position ratio
      const outputScrollTop = outputArea.scrollTop;
      const outputScrollHeight = outputArea.scrollHeight - outputArea.clientHeight;
      const scrollRatio = outputScrollHeight > 0 ? outputScrollTop / outputScrollHeight : 0;

      // Apply same ratio to input
      const inputScrollHeight = inputText.scrollHeight - inputText.clientHeight;
      inputText.scrollTop = scrollRatio * inputScrollHeight;

      setTimeout(() => { isScrollSyncing = false; }, 50);
    }

    convertBtn.addEventListener('click', convertText);
    copyBtn.addEventListener('click', copyToClipboard);
    uploadBtn.addEventListener('click', handleFileUpload);
    fileInput.addEventListener('change', handleFileSelect);
    arrowsBtn.addEventListener('click', extractArrows);
    saveBtn.addEventListener('click', saveAsN4L);
    outputArea.addEventListener('scroll', syncScroll);

    // Allow Ctrl/Cmd + Enter to convert
    inputText.addEventListener('keydown', (e) =>
    {
      if ((e.ctrlKey || e.metaKey) && e.key === 'Enter')
      {
        convertText();
      }
    });

    // Allow Ctrl/Cmd + S to save
    outputArea.addEventListener('keydown', (e) =>
    {
      if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 's')
      {
        e.preventDefault();
        saveAsN4L();
      }
    });

    // Auto-convert on paste
    inputText.addEventListener('paste', () =>
    {
      setTimeout(() =>
      {
        convertText();
      }, 100);
    });

    // Initial focus
    inputText.focus();
  })
  .catch(error =>
  {
    console.error('Error loading app:', error);
    document.getElementById('app').innerHTML = '<p>Error loading application</p>';
  });
