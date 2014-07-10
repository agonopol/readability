package readability

import (
	"strings"
	"testing"
)

func TestRemoveStyle(t *testing.T) {
	body :=
		`
<!DOCTYPE html>
<html lang="en" dir="ltr" class="client-nojs">
<head>
<meta charset="UTF-8" />
<title>The Reluctant Fundamentalist - Wikipedia, the free encyclopedia</title>
<style>a:lang(ar),a:lang(kk-arab),a:lang(mzn),a:lang(ps),a:lang(ur){text-decoration:none}
/* cache key: enwiki:resourceloader:filter:minify-css:7:3904d24a08aa08f6a68dc338f9be277e */</style>
</head>
<body>
	<p>In <i>The Reluctant Fundamentalist</i> the dramatic monologue is conducted by Changez, who starts his conversation addressing to his listener in a street in the centre of Lahore: "Excuse me, sir, but may I be of assistance? Ah, I see I have alarmed you? Do not be frightened by my beard: I am a lover of America. (…) Come, tell me, what were you looking for? Surely, at this time of day, only one thing could have brought you to the district of Old Anarkali —named, as you may be aware, after a courtesan immured for loving a prince—and that is the quest for the perfect cup of tea. Have I guessed correctly? Then allow me, sir, to suggest my favorite among these many establishments."</p>
</body>
</html>
`
	doc, err := Parse([]byte(body))
	if err != nil {
		t.Error(err)
	}
	content, err := doc.Content()
	if err != nil {
		t.Error(err)
	}
	if strings.Contains(content, "style") {
		t.Fail()
	}
}

func TestRemoveScript(t *testing.T) {
	body :=
		`
<!DOCTYPE html>
<html lang="en" dir="ltr" class="client-nojs">
<head>
<script>
if(window.mw){
mw.config.set({"wgCanonicalNamespace":"","wgCanonicalSpecialPageName":false,"wgNamespaceNumber":0,"wgPageName":"The_Reluctant_Fundamentalist","wgTitle":"The Reluctant Fundamentalist","wgCurRevisionId":610311399,"wgRevisionId":610311399,"wgArticleId":7054369,"wgIsArticle":true,"wgIsRedirect":false,"wgAction":"view","wgUserName":null,"wgUserGroups":["*"],"wgCategories":["Use Pakistani English from July 2013","All Wikipedia articles written in Pakistani English","Use dmy dates from July 2013","2007 novels","Ambassador Book Award winning works","Pakistani novels","British novels adapted into films","Frame stories","Novels by Mohsin Hamid","Novels set in Pakistan","Hamish Hamilton books","War in North-West Pakistan fiction","Works originally published in The Paris Review"],"wgBreakFrames":false,"wgPageContentLanguage":"en","wgPageContentModel":"wikitext","wgSeparatorTransformTable":["",""],"wgDigitTransformTable":["",""],"wgDefaultDateFormat":"dmy","wgMonthNames":["","January","February","March","April","May","June","July","August","September","October","November","December"],"wgMonthNamesShort":["","Jan","Feb","Mar","Apr","May","Jun","Jul","Aug","Sep","Oct","Nov","Dec"],"wgRelevantPageName":"The_Reluctant_Fundamentalist","wgIsProbablyEditable":true,"wgRestrictionEdit":[],"wgRestrictionMove":[],"wgWikiEditorEnabledModules":{"toolbar":true,"dialogs":true,"hidesig":true,"preview":false,"previewDialog":false,"publish":false},"wgBetaFeaturesFeatures":[],"wgMediaViewerOnClick":true,"wgVisualEditor":{"isPageWatched":false,"magnifyClipIconURL":"//bits.wikimedia.org/static-1.24wmf11/skins/common/images/magnify-clip.png","pageLanguageCode":"en","pageLanguageDir":"ltr","svgMaxSize":2048,"namespacesWithSubpages":{"6":0,"8":0,"1":true,"2":true,"3":true,"4":true,"5":true,"7":true,"9":true,"10":true,"11":true,"12":true,"13":true,"14":true,"15":true,"100":true,"101":true,"102":true,"103":true,"104":true,"105":true,"106":true,"107":true,"108":true,"109":true,"110":true,"111":true,"447":true,"828":true,"829":true}},"wikilove-recipient":"","wikilove-anon":0,"wgGuidedTourHelpGuiderUrl":"Help:Guided tours/guider","wgFlowTermsOfUseEdit":"By saving changes, you agree to our \u003Ca class=\"external text\" href=\"//wikimediafoundation.org/wiki/Terms_of_use\"\u003ETerms of Use\u003C/a\u003E and agree to irrevocably release your text under the \u003Ca rel=\"nofollow\" class=\"external text\" href=\"//creativecommons.org/licenses/by-sa/3.0\"\u003ECC BY-SA 3.0 License\u003C/a\u003E and \u003Ca class=\"external text\" href=\"//en.wikipedia.org/wiki/Wikipedia:Text_of_the_GNU_Free_Documentation_License\"\u003EGFDL\u003C/a\u003E","wgFlowTermsOfUseSummarize":"By clicking \"Summarize\", you agree to the terms of use for this wiki.","wgFlowTermsOfUseCloseTopic":"By clicking \"Close topic\", you agree to the terms of use for this wiki.","wgFlowTermsOfUseReopenTopic":"By clicking \"Reopen topic\", you agree to the terms of use for this wiki.","wgULSAcceptLanguageList":["en-us","en","ru"],"wgULSCurrentAutonym":"English","wgFlaggedRevsParams":{"tags":{"status":{"levels":1,"quality":2,"pristine":3}}},"wgStableRevisionId":null,"wgCategoryTreePageCategoryOptions":"{\"mode\":0,\"hideprefix\":20,\"showcount\":true,\"namespaces\":false}","wgNoticeProject":"wikipedia","wgWikibaseItemId":"Q597545"});
}
</script>
<script>
</script>
</head>
<body>
	<p>In <i>The Reluctant Fundamentalist</i> the dramatic monologue is conducted by Changez, who starts his conversation addressing to his listener in a street in the centre of Lahore: "Excuse me, sir, but may I be of assistance? Ah, I see I have alarmed you? Do not be frightened by my beard: I am a lover of America. (…) Come, tell me, what were you looking for? Surely, at this time of day, only one thing could have brought you to the district of Old Anarkali —named, as you may be aware, after a courtesan immured for loving a prince—and that is the quest for the perfect cup of tea. Have I guessed correctly? Then allow me, sir, to suggest my favorite among these many establishments."</p>
</body>
</html>
`
	doc, err := Parse([]byte(body))
	if err != nil {
		t.Error(err)
	}
	content, err := doc.Content()
	if err != nil {
		t.Error(err)
	}
	if strings.Contains(content, "script") {
		println(content)
		t.Fail()
	}
}

func TestWalkParagraphs(t *testing.T) {
	body :=
		`
<!DOCTYPE html>
<html lang="en" dir="ltr" class="client-nojs">
<head>
<script>
if(window.mw){
mw.config.set({"wgCanonicalNamespace":"","wgCanonicalSpecialPageName":false,"wgNamespaceNumber":0,"wgPageName":"The_Reluctant_Fundamentalist","wgTitle":"The Reluctant Fundamentalist","wgCurRevisionId":610311399,"wgRevisionId":610311399,"wgArticleId":7054369,"wgIsArticle":true,"wgIsRedirect":false,"wgAction":"view","wgUserName":null,"wgUserGroups":["*"],"wgCategories":["Use Pakistani English from July 2013","All Wikipedia articles written in Pakistani English","Use dmy dates from July 2013","2007 novels","Ambassador Book Award winning works","Pakistani novels","British novels adapted into films","Frame stories","Novels by Mohsin Hamid","Novels set in Pakistan","Hamish Hamilton books","War in North-West Pakistan fiction","Works originally published in The Paris Review"],"wgBreakFrames":false,"wgPageContentLanguage":"en","wgPageContentModel":"wikitext","wgSeparatorTransformTable":["",""],"wgDigitTransformTable":["",""],"wgDefaultDateFormat":"dmy","wgMonthNames":["","January","February","March","April","May","June","July","August","September","October","November","December"],"wgMonthNamesShort":["","Jan","Feb","Mar","Apr","May","Jun","Jul","Aug","Sep","Oct","Nov","Dec"],"wgRelevantPageName":"The_Reluctant_Fundamentalist","wgIsProbablyEditable":true,"wgRestrictionEdit":[],"wgRestrictionMove":[],"wgWikiEditorEnabledModules":{"toolbar":true,"dialogs":true,"hidesig":true,"preview":false,"previewDialog":false,"publish":false},"wgBetaFeaturesFeatures":[],"wgMediaViewerOnClick":true,"wgVisualEditor":{"isPageWatched":false,"magnifyClipIconURL":"//bits.wikimedia.org/static-1.24wmf11/skins/common/images/magnify-clip.png","pageLanguageCode":"en","pageLanguageDir":"ltr","svgMaxSize":2048,"namespacesWithSubpages":{"6":0,"8":0,"1":true,"2":true,"3":true,"4":true,"5":true,"7":true,"9":true,"10":true,"11":true,"12":true,"13":true,"14":true,"15":true,"100":true,"101":true,"102":true,"103":true,"104":true,"105":true,"106":true,"107":true,"108":true,"109":true,"110":true,"111":true,"447":true,"828":true,"829":true}},"wikilove-recipient":"","wikilove-anon":0,"wgGuidedTourHelpGuiderUrl":"Help:Guided tours/guider","wgFlowTermsOfUseEdit":"By saving changes, you agree to our \u003Ca class=\"external text\" href=\"//wikimediafoundation.org/wiki/Terms_of_use\"\u003ETerms of Use\u003C/a\u003E and agree to irrevocably release your text under the \u003Ca rel=\"nofollow\" class=\"external text\" href=\"//creativecommons.org/licenses/by-sa/3.0\"\u003ECC BY-SA 3.0 License\u003C/a\u003E and \u003Ca class=\"external text\" href=\"//en.wikipedia.org/wiki/Wikipedia:Text_of_the_GNU_Free_Documentation_License\"\u003EGFDL\u003C/a\u003E","wgFlowTermsOfUseSummarize":"By clicking \"Summarize\", you agree to the terms of use for this wiki.","wgFlowTermsOfUseCloseTopic":"By clicking \"Close topic\", you agree to the terms of use for this wiki.","wgFlowTermsOfUseReopenTopic":"By clicking \"Reopen topic\", you agree to the terms of use for this wiki.","wgULSAcceptLanguageList":["en-us","en","ru"],"wgULSCurrentAutonym":"English","wgFlaggedRevsParams":{"tags":{"status":{"levels":1,"quality":2,"pristine":3}}},"wgStableRevisionId":null,"wgCategoryTreePageCategoryOptions":"{\"mode\":0,\"hideprefix\":20,\"showcount\":true,\"namespaces\":false}","wgNoticeProject":"wikipedia","wgWikibaseItemId":"Q597545"});
}
</script>
<script>
</script>
</head>
<body>
	<p>
	In <i>The Reluctant Fundamentalist</i> 
	the dramatic monologue is conducted by Changez, who starts his conversation addressing to his listener 
	in a street in the centre of Lahore: "Excuse me, sir, but may I be of assistance? Ah, I see I have alarmed you?
	Do not be frightened by my beard: I am a lover of America. (…) Come, tell me, what were you looking for?
	Surely, at this time of day, only one thing could have brought you to the district of Old Anarkali —named,
	as you may be aware, after a courtesan immured for loving a prince—and that is the quest for the perfect cup
	of tea. Have I guessed correctly? Then allow me, sir, to suggest my favorite among these many establishments."
	</p>
</body>
</html>
`
	doc, err := Parse([]byte(body))
	if err != nil {
		t.Error(err)
	}
	paragraphs, err := doc.doc.Search(`//p|//td`)

	if err != nil {
		t.Error(err)
	}
	if len(paragraphs) == 0 {
		t.Fail()
	}
}
