;(function($){

    function getSelection(el) {
        var start = 0, end = 0, normalizedValue, range,
            textInputRange, len, endRange;

        if (typeof el.selectionStart === 'number' && typeof el.selectionEnd === 'number') {
            start = el.selectionStart;
            end = el.selectionEnd;
        } else {
            range = document.selection.createRange();

            if (range && range.parentElement() === el) {
                len = el.value.length;
                normalizedValue = el.value.replace(/\r\n/g, '\n');

                // Create a working TextRange that lives only in the input
                textInputRange = el.createTextRange();
                textInputRange.moveToBookmark(range.getBookmark());

                // Check if the start and end of the selection are at the very end
                // of the input, since moveStart/moveEnd doesn't return what we want
                // in those cases
                endRange = el.createTextRange();
                endRange.collapse(false);

                if (textInputRange.compareEndPoints('StartToEnd', endRange) > -1) {
                    start = end = len;
                } else {
                    start = -textInputRange.moveStart('character', -len);
                    start += normalizedValue.slice(0, start).split('\n').length - 1;

                    if (textInputRange.compareEndPoints('EndToEnd', endRange) > -1) {
                        end = len;
                    } else {
                        end = -textInputRange.moveEnd('character', -len);
                        end += normalizedValue.slice(0, end).split('\n').length - 1;
                    }
                }
            }
        }

        return {
            start: start,
            end: end
        };
    }

    function setSelectRange(e, start, end) {
        if (!end) end = start;
        if (e.setSelectionRange) {
            // WebKit
            e.focus();
            e.setSelectionRange(start, end);
        } else if (e.createTextRange) {
            // IE
            var range = e.createTextRange();
            range.collapse(true);
            range.moveStart('character', start);
            range.moveEnd('character', end);
            range.select();
        } else if (e.selectionStart) {
            e.selectionStart = start;
            e.selectionEnd = end;
        }
    }

    function skipAhead(s){
        var m = s.match(/^\n+/);
        if(m){
            return m[0].length;
        }
        return 0;
    }

    function skipTrail(s){
        var m = s.match(/\n+$/);
        if(m){
            return m[0].length;
        }
        return 0;
    }

    function NewUndoManager($t, callback){
        var cur = 0;
        var stacks = [];
        var mode = 'none';
        var t = $t.get(0);

        function cbk(){
            if(callback){
                callback(e.canUndo(), e.canRedo());
            }
        }

        function setMode(newMode, repl) {
            if(typeof repl != 'boolean'){
                repl = mode == newMode;
            }
            mode = newMode;
            e.save(repl);
        }

        function getStack(){
            var value = $t.val();
            var stack = $.extend({'value': value}, getSelection(t));
            return stack;
        }

        $t.on('paste drop dragover dragenter', function(){
            setMode('paste', false);
        });

        $t.on('keyup', function(e){
            if (!e.ctrlKey && !e.metaKey) {

                var keyCode = e.keyCode;

                if ((keyCode >= 33 && keyCode <= 40) || (keyCode >= 63232 && keyCode <= 63235)) {
                    // 33 - 40: page up/dn and arrow keys
                    // 63232 - 63235: page up/dn and arrow keys on safari
                    setMode('moving');
                }
                else if (keyCode == 8 || keyCode == 46 || keyCode == 127) {
                    // 8: backspace
                    // 46: delete
                    // 127: delete
                    setMode('deleting');
                }
                else if (keyCode == 13) {
                    // 13: Enter
                    setMode('newlines');
                }
                else if (keyCode == 27) {
                    // 27: escape
                    setMode('escape');
                }
                else if ((keyCode < 16 || keyCode > 20) && keyCode != 91) {
                    // 16-20 are shift, etc.
                    // 91: left window key
                    // I think this might be a little messed up since there are
                    // a lot of nonprinting keys above 20.
                    setMode('typing');
                }
            }
        });

        var e = {
            canRedo: function(){
                return cur < (stacks.length - 1);
            },
            canUndo: function(){
                return cur > 0;
            },
            redo: function(){
                if(e.canRedo()){
                    cur++;
                    var stack = stacks[cur];
                    $t.val(stack.value);
                    setSelectRange(t, stack.start, stack.end);
                }
                cbk();
                $t.focus();
            },
            undo: function(){
                if(e.canUndo()){
                    cur--;
                    var stack = stacks[cur];
                    $t.val(stack.value);
                    setSelectRange(t, stack.start, stack.end);
                }
                cbk();
                $t.focus();
            },
            save: function(repl){
                setTimeout(function(){
                    if(repl){
                        stacks[cur] = getStack();
                    } else if(e.last() !== $t.val()){
                        stacks.push(getStack());
                    }
                    cur = stacks.length - 1;
                    cbk();
                },10);
            },
            last: function(){
                if(stacks.length) {
                    return stacks[stacks.length-1].value;
                }
            }
        };

        stacks.push(getStack());
        cur++;

        return e;
    }

    additionMentions = {}||additionMentions;
    function TextareaComplete($textarea){
        var mentions = [];
        var fetched = false;
        $textarea.textcomplete([
            {
                match: /\B@([\d\w-_]*)$/,
                search: function (term, callback){

                    var cbk = function(){
                        var nums = 0;
                        callback($.map(mentions, function (mention) {
                            if(nums < 5 && mention.indexOf(term) === 0){
                                nums++;
                            }else{
                                mention = null;
                            }
                            return mention;
                        }));
                    };

                    if(fetched){
                        cbk();
                    } else {
                        fetched = true;
                        $.post('/api/user', {action: "get-follows"}, function(d){
                            if(d.success && d.data){
                                $.each(d.data, function(_,d){
                                    additionMentions[d[1]] = d[0];
                                });
                            }
                        }).complete(function(){
                            $.each(additionMentions, function(k,v){
                                var elm = k;
                                if(k != v){
                                    elm = k + ' (' + v + ')';
                                }
                                mentions.push(elm);
                            });
                            cbk();
                        });
                    }
                },
                index: 1,
                replace: function (mention) {
                    var idx = mention.indexOf(' ');
                    if(idx != -1) {
                        mention = mention.substr(0, idx);
                    }
                    return '@' + mention + ' ';
                }
            }
        ])
        // .overlay([
        //     {
        //         match: /\B[@#]([\d\w-_]*)/g,
        //         css: {
        //             'background-color': '#ddd'
        //         }
        //     }
        // ]);
    }

    $(function(){
        $('.markdown-editor').each(function(_,e){
            var $editor = $(e);
            var $state = $editor.find('.md-textarea');
            var $textarea = $state.find('textarea');
            var $preview = $editor.find('.md-preview');
            var $toolbar = $editor.find('.md-toolbar');
            var $undoBtn = $editor.find('[data-meta=undo]');
            var $redoBtn = $editor.find('[data-meta=redo]');
            var url = $editor.data('preview-url');

            var te = $textarea.get(0);

            var saveKey = $editor.data('savekey');
            var cache = "";

            var popup;

            function toggleOtherButtons($btn, disable){
                if(disable){
                    $toolbar.find('.md-btn').not($btn).each(function(_,e){
                        var $e = $(e);
                        $e.data('isDis', $e.hasClass('disabled'));
                        $e.addClass('disabled');
                    });
                } else {
                    $toolbar.find('.md-btn').not($btn).each(function(_,e){
                        var $e = $(e);
                        if(!$e.data('isDis')){
                            $e.removeClass('disabled');
                        }
                    });
                }
            }

            function insertText(v, start, end){
                var sel = getSelection(te);
                var value = $textarea.val();
                var vStart = value.substr(0, sel.start).lastIndexOf('\n') + 1;
                var vv = value.substring(vStart, sel.start);
                if($.trim(vv)){
                    v = '\n' + v;
                    if(start){
                        start += 1;
                    }
                    if(end){
                        end += 1;
                    }
                }
                value = value.substr(0, sel.end) + v + value.substr(sel.end);
                $textarea.val(value);
                setSelectRange(te, start, end);
                undoManager.save();
            }

            var api = {
                 'insertText': insertText,
                 'getSel': function(){
                    return getSelection(te);
                 }
            };

            $editor.data('editor', api);

            if($textarea.val() === '') {
                $textarea.val($.jStorage.get(saveKey));
            }

            var intervalSave = setInterval(function(){
                $.jStorage.set(saveKey, $textarea.val());
            }, 500);

            $textarea.parents('form:first').on('submit', function(){
                clearInterval(intervalSave);
            });

            $textarea.autosize();
            $textarea.css('resize', 'none');
            TextareaComplete($textarea);

            $editor.on('click', '[data-meta=preview]', function(){
                var $e = $(this);
                if($e.hasClass('active')){
                    $state.show();
                    $preview.hide();
                    $e.removeClass('active');
                    toggleOtherButtons($e, false);
                    $textarea.focus();
                } else {
                    $state.hide();
                    $preview.show();
                    $e.addClass('active');
                    toggleOtherButtons($e, true);

                    var n = $.trim($textarea.val());
                    if(n === '') {
                        $preview.html('');
                        return;
                    }
                    if(n == cache) return;

                    cache = n;
                    $.post(url, {'action': 'preview', 'content': n}, function(data){
                        if(data.success){
                            $preview.html(data.preview);
                            if($preview.mdFilter){
                                $preview.mdFilter();
                            }
                        }
                    });
                }
            });

            $textarea.on('keypress', function(e){
                if ((e.ctrlKey || e.metaKey) && (e.keyCode == 89 || e.keyCode == 90)) {
                    e.preventDefault();
                }
            });

            $textarea.on('keydown', function(e){
                var $t = $textarea;
                var te = $t.get(0);
                var st = getSelection(te);
                var start = st.start;
                var end = st.end;
                var value = $t.val();

                var metaKey = e.ctrlKey || e.metaKey;

                switch(e.keyCode){
                case 89:
                    // y
                    if(metaKey){
                        undoManager.redo();
                        e.preventDefault();
                    }
                    break;

                case 90:
                    // z
                    if(metaKey){
                        if(e.shiftKey){
                            undoManager.redo();
                        } else {
                            undoManager.undo();
                        }
                        e.preventDefault();
                    }
                    break;

                case 9:
                    // tab
                    var sel = value.substring(start, end);

                    var df = skipAhead(sel);
                    sel = sel.substr(df);
                    start += df;

                    df = skipTrail(sel);
                    sel = sel.substr(0, sel.length - df);
                    end -= df;

                    if(e.shiftKey){
                        var selStart = value.substr(0, start).lastIndexOf('\n') + 1;
                        var selEnd = value.substr(end).indexOf('\n');
                        if(selEnd == -1) selEnd = 0;
                        selEnd += end;
                        sel = value.substring(selStart, selEnd);
                        var nsel = sel.replace(/^([ ]{1,4}|\t)/gm, '');
                        var dif = sel.length - nsel.length;
                        if(dif > 0){
                            start -= 1;
                            value = value.substr(0, selStart) + nsel + value.substr(selEnd);
                            $t.val(value);
                        }

                        df = skipAhead(value.substr(start));
                        start += df;

                        setSelectRange(te, start, end - dif);
                    } else {
                        sel = '\t' + sel.replace(/\n/g, '\n\t');
                        $t.val(value.substr(0, start) + sel + value.substr(end));
                        setSelectRange(te, start + 1, start + sel.length);
                    }
                    undoManager.save();
                    e.preventDefault();
                    break;
                }
            });

            $editor.find('[data-meta=image]').popover({
                'html': true,
                'container': $editor,
                'title': $editor.find('[rel=image-popover-title]').html(),
                'content': $editor.find('[rel=image-popover-content]').html()
            });

            var onUpload;

            function imageUploadComplete($form, data){
                onUpload = false;
                var $suc = $form.find('.alert-success').hide();
                var $err = $form.find('.alert-danger').hide();
                if(data && data.success){
                    $form.find('[rel=filename]').val('');
                    $form.find('[type=file]').val('');
                    $editor.find('.md-image').find('[name=link]').val(data.link);
                    $suc.show();
                } else if(data && data.msg){
                    $err.text(data.msg);
                    $err.show();
                } else {
                    $err.text($err.data('message'));
                    $err.show();
                }
            }

            $editor.on('click', '[rel=image-insert]', function(){
                var link = $editor.find('.md-image').find('[name=link]').val();
                var sel = getSelection(te);
                insertText("![]("+link+")", sel.start+2);
                $(popup).popover('hide');
            });

            $editor.on('click', '[rel=image-upload]', function(){
                var $form = $editor.find('.md-image-form');
                $form.slideToggle();
                $form.ajaxForm({
                    dataType: 'json',
                    beforeSubmit: function(){
                        if(onUpload || $form.find('[rel=filename]').val() === ''){
                            return false;
                        }
                        onUpload = true;
                    },
                    success: function(data) {
                        imageUploadComplete($form, data);
                    },
                    error: function(){
                        imageUploadComplete($form);
                    }
                });
            });

            $editor.on('click', '[data-meta=code]', function(){
                var sel = getSelection(te);
                if(sel.start != sel.end){
                    var value = $textarea.val();
                    var selv = value.substring(sel.start, sel.end);
                    value = value.substr(0, sel.start) + "\n```go\n"+selv+"\n```";
                    $textarea.val(value);
                    undoManager.save();
                } else {
                    insertText("\n```go\n\n```", sel.start+7);
                }
            });

            var undoManager = NewUndoManager($textarea, function(un, re){
                if(un){
                    $undoBtn.removeClass('disabled');
                } else {
                    $undoBtn.addClass('disabled');
                }
                if(re){
                    $redoBtn.removeClass('disabled');
                } else {
                    $redoBtn.addClass('disabled');
                }
                $textarea.trigger('autosize.resize');
            });

            $undoBtn.on('click', function(){
                undoManager.undo();
            });

            $redoBtn.on('click', function(){
                undoManager.redo();
            });

            $editor.on('show.bs.popover', function(e){
                if(popup&& popup != e.target){
                    $(popup).popover('hide');
                }
                popup = e.target;
            });

            $editor.on('hide.bs.popover', function(e){
                if(popup == e.target){
                    $(popup).data('bs.popover').hoverState = 'out';
                    popup = null;
                }
            });

            $(document).on('mousedown', function(e){
                if(popup && popup != e.target){
                    var $e = $(e.target);
                    var $p = $(popup);
                    if((!$e.hasClass('md-btn') || !($e.hasClass('md-btn') && $e.data('meta') == $p.data('meta'))) &&
                        (!$e.parents('.md-btn').length || !($e.parents('.md-btn').length &&
                        $e.parents('.md-btn:first').data('meta') == $p.data('meta'))) &&
                        !$e.parents('.popover').length){
                        $p.popover('hide');
                    }
                }
            });
        });
    });

})(jQuery);