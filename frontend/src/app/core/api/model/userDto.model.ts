import {BingoBoardDto} from './bingoBoardDto.model';
import {BingoCardDto} from './bingoCardDto.model';

export class UserDto {
    'avatar_url': string;
    'name': string;
    'github_url': string;
    'reddit_url': string;
    'bingo_cards': BingoCardDto[];
    'bingo_boards': BingoBoardDto[];
    'personal_bingo_board': BingoBoardDto;
}
