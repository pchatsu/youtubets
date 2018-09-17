package youtubets_test

const (
	listXML = `
<?xml version="1.0" encoding="utf-8" ?>
<transcript_list docid="-1">
	<track id="0" name="test_en" lang_code="en" lang_original="English" lang_translated="English" lang_default="true"/>
	<track id="1" name="" lang_code="ja" lang_original="日本語" lang_translated="Japanese"/>
</transcript_list>`
	transcriptXML = `
<?xml version="1.0" encoding="utf-8" ?>
<transcript>
	<text start="10.001" dur="6.64">hello</text>
	<text start="20.002" dur="2.136">world.</text>
	<text start="30.345" dur="3.003">&gt;&gt; It&#39;s a test.</text>
</transcript>
`
	includingEmptyLineTranscriptXML = `
<?xml version="1.0" encoding="utf-8" ?>
<transcript>
	<text start="10.001" dur="6.64">hello</text>
	<text start="20.002" dur="2.136">world.</text>
	<text start="30.345" dur="3.003"></text>
</transcript>
`
	emptyListXML = `
<?xml version="1.0" encoding="utf-8" ?>
<transcript_list docid="-1">
</transcript_list>`
)
