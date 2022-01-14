function Hand(props) {
  const { title, cards } = props;
  return (
    <div
      style={{
        borderStyle: 'solid',
        borderRadius: '10px',
        padding: '10px',
        width: '1200px',
        height: '250px'
      }}>
      <p style={{ textAlign: 'left' }}>{title}</p>
      <div style={{ display: 'flex', flexDirection: 'row' }}>
        {cards.map((c) => (
          <p
            style={{
              marginRight: '10px',
              borderStyle: 'solid',
              borderRadius: '10px',
              padding: '2px'
            }}>
            {c.data}
          </p>
        ))}
      </div>
    </div>
  );
}

export default Hand;
