import { useNavigate } from 'react-router-dom'

interface MBTIType {
  code: string
  name: string
  category: string
  color: string
}

const CATEGORY_COLORS = {
  'Sentinels': '#7acae8',
  'Explorers': '#ebcc85',
  'Diplomats': '#80e8be',
  'Analysts': '#c58ed9',
}

export default function MBTIGrid() {
  const navigate = useNavigate()

  const mbtiTypes: MBTIType[][] = [
    // Row 1: Sentinels (SJ)
    [
      { code: 'ISTJ', name: 'The Logistician', category: 'Sentinels', color: CATEGORY_COLORS.Sentinels },
      { code: 'ISFJ', name: 'The Defender', category: 'Sentinels', color: CATEGORY_COLORS.Sentinels },
      { code: 'ESTJ', name: 'The Executive', category: 'Sentinels', color: CATEGORY_COLORS.Sentinels },
      { code: 'ESFJ', name: 'The Consul', category: 'Sentinels', color: CATEGORY_COLORS.Sentinels },
    ],
    // Row 2: Explorers (SP)
    [
      { code: 'ISTP', name: 'The Virtuoso', category: 'Explorers', color: CATEGORY_COLORS.Explorers },
      { code: 'ISFP', name: 'The Adventurer', category: 'Explorers', color: CATEGORY_COLORS.Explorers },
      { code: 'ESTP', name: 'The Entrepreneur', category: 'Explorers', color: CATEGORY_COLORS.Explorers },
      { code: 'ESFP', name: 'The Entertainer', category: 'Explorers', color: CATEGORY_COLORS.Explorers },
    ],
    // Row 3: Diplomats (NF)
    [
      { code: 'INFJ', name: 'The Advocate', category: 'Diplomats', color: CATEGORY_COLORS.Diplomats },
      { code: 'INFP', name: 'The Mediator', category: 'Diplomats', color: CATEGORY_COLORS.Diplomats },
      { code: 'ENFJ', name: 'The Protagonist', category: 'Diplomats', color: CATEGORY_COLORS.Diplomats },
      { code: 'ENFP', name: 'The Campaigner', category: 'Diplomats', color: CATEGORY_COLORS.Diplomats },
    ],
    // Row 4: Analysts (NT)
    [
      { code: 'INTJ', name: 'The Architect', category: 'Analysts', color: CATEGORY_COLORS.Analysts },
      { code: 'INTP', name: 'The Logician', category: 'Analysts', color: CATEGORY_COLORS.Analysts },
      { code: 'ENTJ', name: 'The Commander', category: 'Analysts', color: CATEGORY_COLORS.Analysts },
      { code: 'ENTP', name: 'The Debater', category: 'Analysts', color: CATEGORY_COLORS.Analysts },
    ],
  ]

  const handleCardClick = (mbtiCode: string) => {
    navigate(`/mbti/${mbtiCode}`)
  }

  return (
    <div className="mbti-grid-container">
      <div className="page-header">
        <div className="header-title">MBTI PERSONALITY TYPE CLASSIFICATION</div>
        <div className="header-subtitle">SELECT A TYPE TO EXPLORE</div>
      </div>

      {mbtiTypes.map((row, rowIndex) => (
        <div
          key={rowIndex}
          className="mbti-category-section"
          style={{
            '--card-color': row[0].color,
          } as React.CSSProperties}
        >
          <div className="mbti-category-header">
            ═══ {row[0].category.toUpperCase()} ═══
          </div>
          <div className="mbti-row">
            {row.map((type) => (
              <div
                key={type.code}
                className="mbti-card"
                onClick={() => handleCardClick(type.code)}
                style={{
                  '--card-color': type.color,
                } as React.CSSProperties}
              >
                <div className="mbti-card-code">{type.code}</div>
                <div className="mbti-card-name">{type.name}</div>
                <div className="mbti-card-hint">[ CLICK TO VIEW AGENTS ]</div>
              </div>
            ))}
          </div>
        </div>
      ))}

      <div className="mbti-footer">
        <div className="mbti-info">
          ℹ Click any personality type to view agents with that MBTI classification
        </div>
      </div>
    </div>
  )
}
