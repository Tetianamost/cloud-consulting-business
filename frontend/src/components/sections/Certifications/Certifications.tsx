import React from 'react';
import styled from 'styled-components';
import { motion } from 'framer-motion';
import { theme } from '../../../styles/theme';
import Section from '../../ui/section';
import CertificationBadge from './CertificationBadge';
import { FaAward, FaCertificate } from 'react-icons/fa';
import { SiAmazon } from 'react-icons/si';

const certifications = [
	{
		id: 5,
		title: 'AWS Advanced Networking – Specialty',
		description:
			'Expert-level knowledge in designing and implementing complex AWS network architectures',
		image: SiAmazon,
		color: '#FDCB6E',
	},
	{
		id: 3,
		title: 'AWS Certified Database – Specialty',
		description:
			'Advanced skills in designing, implementing, and managing AWS database solutions',
		image: FaAward,
		color: '#4ECDC4',
	},
	{
		id: 1,
		title: 'AWS Certified Solutions Architect – Associate',
		description:
			'Validated expertise in designing distributed systems and AWS cloud architectures',
		image: SiAmazon,
		color: '#FF9900',
	},
	{
		id: 2,
		title: 'AWS Certified Developer – Associate',
		description:
			'Certified in developing and maintaining AWS-based applications and services',
		image: FaCertificate,
		color: '#FF6B6B',
	},
	{
		id: 4,
		title: 'AWS Certified Cloud Practitioner',
		description:
			'Comprehensive understanding of AWS cloud fundamentals and best practices',
		image: FaCertificate,
		color: '#45B7D1',
	},
];

const SectionTitle = styled(motion.h2)`
	text-align: center;
	margin-bottom: ${theme.space[3]};
	color: ${theme.colors.white};
`;

const SectionSubtitle = styled(motion.p)`
	text-align: center;
	max-width: 600px;
	margin: 0 auto ${theme.space[8]};
	color: ${theme.colors.gray300};
	font-size: ${theme.fontSizes.lg};
`;

const CertificationsGrid = styled.div`
	display: grid;
	grid-template-columns: repeat(3, 1fr);
	gap: ${theme.space[6]};

	@media (max-width: ${theme.breakpoints.lg}) {
		grid-template-columns: repeat(2, 1fr);
	}

	@media (max-width: ${theme.breakpoints.sm}) {
		grid-template-columns: 1fr;
	}
`;

// Animation variants
const titleVariants = {
	hidden: { opacity: 0, y: 20 },
	visible: {
		opacity: 1,
		y: 0,
		transition: {
			duration: 0.5,
		},
	},
};

const Certifications: React.FC = () => {
	return (
		<Section id="certifications" background="primary">
			<SectionTitle
				initial="hidden"
				whileInView="visible"
				viewport={{ once: true }}
				variants={titleVariants}
			>
				Our AWS Certifications
			</SectionTitle>
			<SectionSubtitle
				initial="hidden"
				whileInView="visible"
				viewport={{ once: true }}
				variants={titleVariants}
			>
				We bring our professional AWS cloud expertise to deliver high quality
				consulting services
			</SectionSubtitle>

			<CertificationsGrid>
				{certifications.map((certification, index) => (
					<CertificationBadge
						key={certification.id}
						certification={certification}
						index={index}
					/>
				))}
			</CertificationsGrid>
		</Section>
	);
};

const MemoizedCertifications = React.memo(Certifications);
export default MemoizedCertifications;