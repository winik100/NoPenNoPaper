const checkbox_ids = ['AnthropologieEdit', 'ArchäologieEdit', 'AutofahrenEdit', 'BibliotheksnutzungEdit', 'BuchführungEdit', 'CharmeEdit', 'EinschüchternEdit', 
    'Elektrische ReparaturenEdit', 'Erste HilfeEdit', 'FinanzkraftEdit', 'GeschichteEdit', 'HorchenEdit', 'KaschierenEdit', 'KletternEdit', 'Mechanische ReparaturenEdit', 'MedizinEdit',
    'NaturkundeEdit', 'OkkultismusEdit', 'OrientierungEdit', 'PsychoanalyseEdit', 'PsychologieEdit', 'RechtswesenEdit', 'ReitenEdit', 'SchließtechnikEdit', 'Schweres GerätEdit', 
    'SchwimmenEdit', 'SpringenEdit', 'SpurensucheEdit', 'ÜberredenEdit', 'ÜberzeugenEdit', 'Verborgen bleibenEdit', 'Verborgenes erkennenEdit', 'VerkleidenEdit', 'WerfenEdit', 'Werte schätzenEdit'];

for (const id of checkbox_ids){
    var el = document.getElementById(id);
    if (el) {
        el.addEventListener('click', function() {
            newId = id.substring(0, id.length - 4)
            elem = document.getElementById(newId)
            elem.disabled = !this.checked

            newId = id.substring(0, id.length - 4).concat('Val')
            elem = document.getElementById(newId)
            elem.disabled = !this.checked
        })
    }
}